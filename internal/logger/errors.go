package logger

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

// Err returns an "error" attribute that is safe to log. MySQL server errors
// can echo column values (e.g. "Duplicate entry '<token>' for key 'token'"),
// so for those only the error number is logged.
func Err(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "")
	}
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return slog.String("error", "mysql error "+strconv.Itoa(int(me.Number)))
	}
	return slog.String("error", err.Error())
}
