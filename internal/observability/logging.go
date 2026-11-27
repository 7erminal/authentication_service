package observability

import (
	"io"
	"log"
	"log/slog"
	"strings"
	"sync"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

var registerLogger sync.Once

// Access events are always enabled, independently of application severity.
var accessLogger *slog.Logger

// ConfigureLogging intentionally emits JSON to the supplied sink (stdout in
// production). File collection/rotation belongs to the deployment log collector.
func ConfigureLogging(out io.Writer, configuredLevel string) error {
	var level slog.Level
	if configuredLevel == "" {
		configuredLevel = "INFO"
	}
	if err := level.UnmarshalText([]byte(strings.ToUpper(configuredLevel))); err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})))
	accessLogger = slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	orm.Debug = false
	orm.DebugLog.SetOutput(io.Discard)
	// Framework messages have no request context and may contain URLs, SQL or
	// secrets. Do not forward their free-form content. Application request errors
	// are logged through trace.Logger instead.
	registerLogger.Do(func() { logs.Register("safe_slog", func() logs.Logger { return &frameworkLogger{} }) })
	logs.Reset()
	if err := logs.SetLogger("safe_slog"); err != nil {
		return err
	}
	logs.SetLevel(logs.LevelWarning) // also disables Beego's dev-mode access line
	log.SetFlags(0)
	log.SetOutput(io.Discard)
	return nil
}

type frameworkLogger struct{}

func (*frameworkLogger) Init(string) error              { return nil }
func (*frameworkLogger) Destroy()                       {}
func (*frameworkLogger) Flush()                         {}
func (*frameworkLogger) SetFormatter(logs.LogFormatter) {}

// Beego's global logger cannot carry request context. Suppress these diagnostics
// rather than emit uncorrelated or potentially sensitive free-form messages.
// Access, application errors, startup failures and panics have dedicated events.
func (*frameworkLogger) WriteMsg(*logs.LogMsg) error { return nil }
