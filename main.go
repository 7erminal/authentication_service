package main

import (
	_ "authentication_service/routers"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"

	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	sqlConn, err := beego.AppConfig.String("sqlconn")
	if err != nil {
		logs.Error("%s", err)
	}
	orm.RegisterDataBase("default", "mysql", sqlConn)
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/application.log"}`)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8000", "http://152.67.134.169", "http://13.40.60.131:8001", "http://167.86.115.44:8002"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
