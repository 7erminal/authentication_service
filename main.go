package main

import (
	_ "authentication_service/routers"
	"fmt"
	"net/url"
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"

	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/lib/pq"
)

func setup() {
	fmt.Println("Setting up the database connection...")
	// 1. Fetch values from app.conf
	// val, _ := beego.AppConfig.GetSection("default")
	user, err := beego.AppConfig.String("pguser")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	pass, err := beego.AppConfig.String("pgpass")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	host, err := beego.AppConfig.String("pghost")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	port, err := beego.AppConfig.String("pgport")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	dbName, err := beego.AppConfig.String("pgdb")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	pgTest, err := beego.AppConfig.String("pgtest")
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch configuration: %v", err))
	}
	logs.Info("Postgres test value is: ", pgTest)

	// 2. Register the driver name
	orm.RegisterDriver("postgres", orm.DRPostgres)

	if user == "" {
		user = os.Getenv("DB_USER")
	}
	if pass == "" {
		pass = os.Getenv("DB_PASSWORD")
	}
	if host == "" {
		host = os.Getenv("DB_HOST")
	}
	if port == "" {
		port = os.Getenv("DB_PORT")
	}
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}

	fmt.Printf("Database credentials are: user=%s, host=%s, port=%s, dbName=%s\n", user, host, port, dbName)

	// 3. Construct connection string (handle special characters in password)
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		user, url.QueryEscape(pass), host, port, dbName,
	)

	// 4. Register the default database alias
	// Max connection pool limits can be configured here as the 4th/5th optional arguments
	dberr := orm.RegisterDataBase("default", "postgres", dataSource)
	if dberr != nil {
		panic(fmt.Sprintf("Failed to register database: %v", dberr))
	}
}

func main() {
	// sqlConn, err := beego.AppConfig.String("sqlconn")
	// if err != nil {
	// 	logs.Error("%s", err)
	// }
	// orm.RegisterDataBase("default", "mysql", sqlConn)
	pgTest, err := beego.AppConfig.String("pgtest")
	if err != nil {
		logs.Error("%s", err)
	}
	logs.Info("Postgres test ", pgTest)

	custUrl, err := beego.AppConfig.String("customerBaseUrl")
	if err != nil {
		logs.Error("%s", err)
	}
	logs.Info("Customer Base URL is ", custUrl)
	setup()

	logs.Info("Starting Authentication Service...")

	logs.SetLogger(logs.AdapterFile, `{"filename":"../logs/authentication_application.log"}`)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000", "http://localhost:8000", "http://152.67.134.169", "http://13.40.60.131:8001", "http://167.86.115.44:8002", "http://5.252.55.191", "http://185.249.227.127:9001", "http://185.249.227.127:9002", "makufoodsltd.net", "https://makufoodsltd.net", "https://www.makufoodsltd.com", "https://makufoodsltd.com", "makufoodsltd.com", "https://admin.bridgeafrica.group", "https://mestechgh.com", "https://admin.mestechgh.com", "https://client.mestechgh.com", "https://authentication.mestechgh.com"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	logs.Info("Initializing server...")

	orm.Debug = true
	beego.Run()
}
