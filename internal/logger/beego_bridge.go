package logger

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"

	"github.com/beego/beego/v2/core/logs"
)

// The Beego bridge sends legacy logs.Info/Error/... calls (not yet migrated
// to slog) to the JSON logger so the whole service emits one format. These
// calls have no request context, so their ID fields are empty; they are
// tagged logger="beego" to make them easy to find and migrate.

const bridgeAdapter = "slog_bridge"

var (
	bridgeOnce   sync.Once
	bridgeMu     sync.RWMutex
	bridgeTarget *slog.Logger
)

// devRouterLine matches the per-request line Beego prints in dev run mode
// (router.go: "|   127.0.0.1| 200 |  1.2ms|   match| POST  /v1/..."). We emit
// our own access log, so this duplicate is dropped.
var devRouterLine = regexp.MustCompile(`^\|\s*[^|]*\|.*\d{3}.*\|\s*[^|]*\|\s*(match|nomatch)\|`)

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

type bridge struct{}

func (bridge) Init(string) error              { return nil }
func (bridge) Destroy()                       {}
func (bridge) Flush()                         {}
func (bridge) SetFormatter(logs.LogFormatter) {}
func (bridge) WriteMsg(lm *logs.LogMsg) error { return writeBeego(lm) }
func newBridge() logs.Logger                  { return bridge{} }

func installBeegoBridge(l *slog.Logger) {
	bridgeMu.Lock()
	bridgeTarget = l
	bridgeMu.Unlock()
	bridgeOnce.Do(func() { logs.Register(bridgeAdapter, newBridge) })
	// Replace all adapters (console/file) with the bridge. The slog handler
	// does the level filtering, so let every Beego level through.
	logs.Reset()
	_ = logs.SetLogger(bridgeAdapter)
	logs.SetLevel(logs.LevelDebug)
}

func writeBeego(lm *logs.LogMsg) error {
	msg := lm.Msg
	if len(lm.Args) > 0 {
		msg = fmt.Sprintf(lm.Msg, lm.Args...)
	}
	msg = strings.TrimSpace(ansi.ReplaceAllString(msg, ""))
	if devRouterLine.MatchString(msg) {
		return nil
	}

	bridgeMu.RLock()
	l := bridgeTarget
	bridgeMu.RUnlock()
	if l == nil {
		return nil
	}
	l.Log(context.Background(), beegoLevel(lm.Level), msg, slog.String("logger", "beego"))
	return nil
}

func beegoLevel(level int) slog.Level {
	switch {
	case level <= logs.LevelError: // emergency, alert, critical, error
		return slog.LevelError
	case level == logs.LevelWarn:
		return slog.LevelWarn
	case level == logs.LevelDebug:
		return slog.LevelDebug
	default: // notice, info
		return slog.LevelInfo
	}
}
