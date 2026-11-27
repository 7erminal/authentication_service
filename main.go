package main

import (
	"authentication_service/internal/observability"
	"authentication_service/internal/trace"
	_ "authentication_service/routers"
	"context"
	"os"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"

	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level, _ = beego.AppConfig.String("loglevel")
	}
	if err := observability.ConfigureLogging(os.Stdout, level); err != nil {
		// Configuration failures must not expose arbitrary env/config values.
		_ = observability.ConfigureLogging(os.Stdout, "INFO")
		trace.Logger(context.Background()).Error("invalid logging level configuration")
		os.Exit(1)
	}
	observability.Install()
	sqlConn, err := beego.AppConfig.String("sqlconn")
	if err != nil {
		trace.Logger(context.Background()).Error("database configuration unavailable")
		os.Exit(1)
	}
	orm.Debug = false
	if err := orm.RegisterDataBase("default", "mysql", sqlConn); err != nil {
		trace.Logger(context.Background()).Error("database registration failed")
		os.Exit(1)
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000", "http://localhost:8000", "http://152.67.134.169", "http://13.40.60.131:8001", "http://167.86.115.44:8002", "http://5.252.55.191", "http://185.249.227.127:9001", "http://185.249.227.127:9002", "makufoodsltd.net", "https://makufoodsltd.net", "https://www.makufoodsltd.com", "https://makufoodsltd.com", "makufoodsltd.com", "https://admin.bridgeafrica.group", "https://mestechgh.com", "https://admin.mestechgh.com", "https://client.mestechgh.com", "https://authentication.mestechgh.com"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", trace.RequestHeader, trace.TraceHeader, trace.SessionHeader},
		ExposeHeaders:    []string{"Content-Length", trace.RequestHeader, trace.TraceHeader},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	trace.Logger(context.Background()).Info("service starting", "http_port", beego.BConfig.Listen.HTTPPort)
	beego.RunWithMiddleWares("", observability.Middleware)
}
