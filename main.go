package main

import (
	"log/slog"
	"os"

	"authentication_service/controllers"
	"authentication_service/filters"
	"authentication_service/internal/logger"
	"authentication_service/internal/trace"
	_ "authentication_service/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"

	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	// Structured JSON logs on stdout (captured by pm2). Level: LOG_LEVEL env,
	// else "loglevel" in app.conf, else info. Legacy logs.* calls are bridged.
	logLevel, _ := beego.AppConfig.String("loglevel")
	logger.Init(os.Stdout, logger.ResolveLevel(logLevel))

	sqlConn, err := beego.AppConfig.String("sqlconn")
	if err != nil {
		slog.Error("sqlconn not configured", "error", err.Error())
	}
	orm.RegisterDataBase("default", "mysql", sqlConn)

	// Tracing, access log and panic recovery must be installed before CORS so
	// preflight and every other response carries X-Request-Id / X-Trace-Id.
	filters.Install(beego.BeeApp)
	beego.ErrorController(&controllers.ErrorController{})

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000", "http://localhost:8000", "http://152.67.134.169", "http://13.40.60.131:8001", "http://167.86.115.44:8002", "http://5.252.55.191", "http://185.249.227.127:9001", "http://185.249.227.127:9002", "makufoodsltd.net", "https://makufoodsltd.net", "https://www.makufoodsltd.com", "https://makufoodsltd.com", "makufoodsltd.com", "https://admin.bridgeafrica.group", "https://mestechgh.com", "https://admin.mestechgh.com", "https://client.mestechgh.com", "https://authentication.mestechgh.com"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", trace.HeaderRequestID, trace.HeaderTraceID, trace.HeaderSessionID},
		ExposeHeaders:    []string{"Content-Length", trace.HeaderRequestID, trace.HeaderTraceID},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// SQL debug logging prints full statements including password hashes and
	// tokens, so it is opt-in only.
	orm.Debug = os.Getenv("ORM_DEBUG") == "true"

	if !controllers.GoogleOAuthConfigured() {
		slog.Warn("google oauth disabled: credentials not set",
			"env", controllers.EnvGoogleClientID+"/"+controllers.EnvGoogleClientSecret)
	}
	beego.Run()
}
