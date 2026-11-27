package models

import (
	"authentication_service/internal/observability"
	"authentication_service/internal/trace"
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/beego/beego/v2/client/orm"
)

// A standard-library SQL driver verifies context at the actual driver boundary,
// without a live database or a mocking dependency.
type contextDriver struct{ contexts chan context.Context }

func (d *contextDriver) Open(string) (driver.Conn, error) { return &contextConnection{d}, nil }

type contextConnection struct{ driver *contextDriver }

func (*contextConnection) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected context-free prepare")
}
func (*contextConnection) Close() error { return nil }
func (*contextConnection) Begin() (driver.Tx, error) {
	return nil, errors.New("unexpected transaction")
}
func (*contextConnection) Ping(context.Context) error { return nil }
func (c *contextConnection) ExecContext(ctx context.Context, _ string, _ []driver.NamedValue) (driver.Result, error) {
	c.driver.contexts <- ctx
	return contextResult{}, nil
}
func (c *contextConnection) QueryContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "TIMEDIFF") {
		return &contextRows{columns: []string{"tz"}, values: []driver.Value{"00:00:00"}}, nil
	}
	if strings.Contains(query, "information_schema") {
		return &contextRows{columns: []string{"engine", "transactions"}, values: []driver.Value{"INNODB", true}}, nil
	}
	c.driver.contexts <- ctx
	return &contextRows{columns: []string{"id"}}, nil
}

type contextResult struct{}

func (contextResult) LastInsertId() (int64, error) { return 42, nil }
func (contextResult) RowsAffected() (int64, error) { return 1, nil }

type contextRows struct {
	columns []string
	values  []driver.Value
}

func (r *contextRows) Columns() []string { return r.columns }
func (*contextRows) Close() error        { return nil }
func (r *contextRows) Next(dest []driver.Value) error {
	if r.values == nil {
		return io.EOF
	}
	copy(dest, r.values)
	r.values = nil
	return nil
}

var driverOnce sync.Once
var modelTestDriver = &contextDriver{contexts: make(chan context.Context, 8)}

func TestModelContextReachesSQL(t *testing.T) {
	var output bytes.Buffer
	_ = observability.ConfigureLogging(&output, "DEBUG")

	d := modelTestDriver
	driverOnce.Do(func() {
		sql.Register("trace_context_test", d)
		if err := orm.RegisterDriver("trace_context_test", orm.DRMySQL); err != nil {
			t.Fatal(err)
		}
		db, err := sql.Open("trace_context_test", "")
		if err != nil {
			t.Fatal(err)
		}
		if err := orm.AddAliasWthDB("default", "trace_context_test", db); err != nil {
			t.Fatal(err)
		}
	})
	ids := trace.IDs{SessionID: "session_123", TraceID: "trace_123", RequestID: "request_123"}
	ctx := trace.WithIDs(context.Background(), ids)
	id, err := AddAccessTokensContext(ctx, &AccessTokens{Token: "TOKEN_SECRET"})
	if err != nil || id != 42 {
		t.Fatalf("insert: %d %v", id, err)
	}
	if got := <-d.contexts; got != ctx || trace.FromContext(got) != ids {
		t.Fatal("insert lost request context")
	}
	_, err = GetAccessTokensByTokenContext(ctx, "TOKEN_SECRET")
	if !errors.Is(err, orm.ErrNoRows) {
		t.Fatalf("lookup behavior changed: %v", err)
	}
	if got := <-d.contexts; got != ctx {
		t.Fatal("query lost request context")
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	_, err = AddAccessTokensContext(cancelled, &AccessTokens{Token: "TOKEN_SECRET"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("DB cancellation lost: %v", err)
	}
	if len(d.contexts) != 0 {
		t.Fatal("canceled query reached driver")
	}
	if !strings.Contains(output.String(), `"record_id":42`) || !strings.Contains(output.String(), `"request_id":"request_123"`) {
		t.Fatal("record ID and trace not logged together")
	}
	if strings.Contains(output.String(), "TOKEN_SECRET") {
		t.Fatal("SQL parameters leaked")
	}
}
