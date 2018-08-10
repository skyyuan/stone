package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"stone/common"
	"stone/common/auth"
	"stone/common/webcomm"
	"stone/locale"
	"stone/nsqs"

	"stone/cmd/web/api/v1"
	"stone/middleware"

	"github.com/alecthomas/kingpin"
	"github.com/labstack/echo"
	emw "github.com/labstack/echo/middleware"

)

var (
	// Version version
	Version = "0.0.1"
	app     = kingpin.New("app", "Web application server.").DefaultEnvars()
	port    = app.Flag("port", "Server port for listening.").Short('p').Default("8080").String()
	//port    = app.Flag("port", "Server port for listening.").Short('p').Default("8080").String()
	authoff = app.Flag("authoff", "Turnoff application authrization").Default("false").Bool()

	cmdRun = app.Command("run", "Run application").Default()

	cmdGenAppKey  = app.Command("genappkey", "Generate an application key.")
	genAppKeyName = cmdGenAppKey.Flag("name", "App name").Required().String()
	genAppKeyDesc = cmdGenAppKey.Flag("desc", "App description").String()
)

func main() {
	conf := configure()
	kingpin.Version(Version)

	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	// common init
	common.InitDB(conf)
	defer common.DBClose()
	dDbMigrate()

	switch parsedCmd {
	// Generate app key
	case cmdGenAppKey.FullCommand():
		appAuth, _ := auth.GenerateAppAuth(*genAppKeyName, *genAppKeyDesc)
		fmt.Printf("generated [%s] authorization: key [%s]  secret [%s]\n", appAuth.AppName, appAuth.AppKey, appAuth.Secret)
		return
	}

	err := nsqs.InitConfig(conf.Nsqconfig)
	if err != nil {
		common.Logger.Fatal(err)
	}
	e := echo.New()
	common.EchoInit(e, conf)
	locale.Init()

	// middlewares
	e.Pre(emw.RemoveTrailingSlash())
	e.Pre(middleware.NoCache())
	e.Pre(middleware.Heartbeat("/ping"))
	e.Use(middleware.RequestID())
	e.Use(emw.Secure())
	reqlogger := common.RequestLog(conf.Reqlog())
	e.Use(middleware.Logger(reqlogger), middleware.Recover())
	// 鉴权
	//if !*authoff {
	//	e.Use(middleware.AppAuth())
	//}

	// actions
	v1.RegisterAPI(e)

	srvAddr := ":" + *port

	e.Logger.Infof("Listening and serving HTTP on %s\n", srvAddr)
	// Start server
	go func() {
		if err := e.Start(srvAddr); err != nil {
			common.Logger.Info(err)
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("Server exist")
}

func configure() *webcomm.Config {
	conf := &webcomm.Config{}
	app.Flag("db", "Database connection URL, only support mysql.").
		//PlaceHolder("USER:PWD@tcp(DBURL:DBPORT)/DBSCHEMA?charset=utf8&parseTime=True&loc=Local").
		PlaceHolder("user:"+""+"@tcp(127.0.0.1:3306)/DBSCHEMA?charset=utf8mb4&parseTime=True&loc=Local").
		StringVar(&conf.CmMysqlURL)
	app.Flag("dbidle", "Database idel connection numbers.").
		Default("10").
		IntVar(&conf.CmMysqlIdle)
	app.Flag("dbmax", "Database max-open.").
		Default("100").
		IntVar(&conf.CmMysqlMaxOpen)
	app.Flag("reqlog", "Request log, support file only").
		PlaceHolder("/tmp/wallet/eth-req.log").
		StringVar(&conf.CmReqlog)
	app.Flag("debug", "Enable debug mode").Default("false").
		BoolVar(&conf.CmDebug)
	conf.Nsqconfig = &nsqs.SimpleConfig{}
	app.Flag("nsqd", "Nsqd node address").Default("127.0.0.1:4151").
		StringVar(&conf.Nsqconfig.NsqAddress)
	app.Flag("nsqmif", "Nsqd max-in-flight").Default("100").
		IntVar(&conf.Nsqconfig.MaxInFlight)
	app.Flag("nsqlookups", "Nsq lookups address").Default("127.0.0.1:4161").
		StringsVar(&conf.Nsqconfig.Lookups)
	return conf
}