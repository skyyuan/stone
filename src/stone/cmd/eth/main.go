package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"stone/cmd/eth/api/v2"

	"stone/common"
	"stone/common/ethcomm"
	"stone/locale"
	"stone/service/eth"

	"stone/cmd/eth/api/v1"
	"stone/middleware"

	"github.com/alecthomas/kingpin"
	"github.com/labstack/echo"
	emw "github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

var (
	// Version version
	Version = "0.0.1"
	app     = kingpin.New("app", "Eth applicaton server").DefaultEnvars()
	cmdRun  = app.Command("run", "Run application").Default()
	port    = app.Flag("port", "Server port for listening.").Short('p').Default("8080").String()
	authoff = app.Flag("authoff", "Turnoff application authrization").Default("false").Bool()
)

func main() {
	conf := configure()
	kingpin.Version(Version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	e := echo.New()
	common.EchoInit(e, conf)
	e.Validator = &common.SimpleValidator{Validator: validator.New()}
	common.InitDB(conf)
	defer common.DBClose()
	ethcomm.InitDB(conf)
	defer ethcomm.DBClose()
	eth.ServiceInit(conf)
	locale.Init()
	dbMigrate()

	// middlewares
	e.Pre(emw.RemoveTrailingSlash())
	e.Pre(middleware.NoCache())
	e.Pre(middleware.Heartbeat("/ping"))
	e.Use(middleware.RequestID())
	e.Use(emw.Secure())
	reqlogger := common.RequestLog(conf.Reqlog())
	e.Use(middleware.Logger(reqlogger), middleware.Recover())
	// 鉴权
	if !*authoff {
		e.Use(middleware.AppAuth())
	}
	e.Use(emw.CORSWithConfig(emw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// actions
	v1.RegisterAPI(e)
	v2.RegisterAPI(e)

	srvAddr := ":" + *port

	e.Logger.Infof("Listening and serving HTTP on %s\n", srvAddr)
	// Start server
	go func() {
		if err := e.Start(srvAddr); err != nil {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	eth.ServiceDone()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("Server exist")
}

func configure() *ethcomm.EthConfig {
	conf := &ethcomm.EthConfig{}
	app.Flag("db", "Database connection URL, only support mysql.").
		PlaceHolder("USER:PWD@tcp(DBURL:DBPORT)/DBSCHEMA?charset=utf8&parseTime=True&loc=Local").
		Required().
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
	app.Flag("gethendpoint", "Geth endpoint").
		PlaceHolder("http://IP:PORT").
		Required().
		StringVar(&conf.EthComGethEndpoint)
	app.Flag("ethdb", "Eth database connection URL, only support mysql.").
		PlaceHolder("USER:PWD@tcp(DBURL:DBPORT)/DBSCHEMA?charset=utf8&parseTime=True&loc=Local").
		Required().
		StringVar(&conf.EthComMysqlURL)
	app.Flag("ethdbidle", "Eth database idel connection numbers.").
		Default("10").
		IntVar(&conf.EthComMysqlIdle)
	app.Flag("ethdbmax", "Eth database max-open.").
		Default("100").
		IntVar(&conf.EthComMysqlMaxOpen)
	app.Flag("testnet", "Eth testnet.").
		Default("false").
		BoolVar(&conf.EthComTestNet)

	return conf
}
