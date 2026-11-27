package filters

import (
	beego "github.com/beego/beego/v2/server/web"
)

// Install wires tracing, access logging and panic recovery into app. Call it
// before registering other filters (e.g. CORS) so the tracing filter runs
// first and every response carries the IDs.
func Install(app *beego.HttpServer) {
	app.Cfg.RecoverFunc = RecoverPanic
	// Our access log replaces Beego's (LogAccess / dev-mode router line).
	app.Cfg.Log.AccessLogs = false
	// Outermost wrapper: one access-log line per request, panic safety net.
	app.InsertFilterChain("*", AccessLog)
	// BeforeRouter tracing filter (idempotent; see Tracing).
	app.InsertFilter("*", beego.BeforeRouter, Tracing)
}
