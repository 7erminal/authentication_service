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
	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/authentication_application.log"}`)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "https://admin.mestechltd.com", "https://client.mestechltd.com", "https://authentication.mestechltd.com"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	orm.Debug = true
	beego.Run()
}
