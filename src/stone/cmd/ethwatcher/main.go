package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron"
	"stone/common"
	"stone/common/ethcomm"
	"stone/nsqs"
	"stone/service/eth"
	"stone/service/eth/syncdb"
)

var (
	// Version version
	Version = "0.0.1"
	app     = kingpin.New("app", "Eth watcher server").DefaultEnvars()
	cmdRun  = app.Command("run", "Run application").Default()
)

func main() {
	conf := configure()
	kingpin.Version(Version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	eth.ServiceInit(conf)
	ethcomm.InitDB(conf)
	defer ethcomm.DBClose()
	dbMigrate()
	common.Logger = log.New("ethwatcher")

	err := nsqs.InitConfig(conf.Nsqconfig)
	if err != nil {
		common.Logger.Fatal(err)
	}
	syncdb.DisableAutoErc20 = conf.DisableAutoErc20
	syncdb.CheckAndRepairBlockData()

	c := cron.New()
	c.AddFunc("@every 10s", syncdb.SyncEthDB)
	c.Start()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	eth.ServiceDone()
	c.Stop()
}

func configure() *ethcomm.EthConfig {
	conf := &ethcomm.EthConfig{}
	app.Flag("ethdb", "Database connection URL, only support mysql.").
		PlaceHolder("USER:PWD@tcp(DBURL:DBPORT)/DBSCHEMA?charset=utf8&parseTime=True&loc=Local").
		Required().
		StringVar(&conf.EthComMysqlURL)
	app.Flag("ethdbidle", "Database idel connection numbers.").
		Default("20").
		IntVar(&conf.EthComMysqlIdle)
	app.Flag("ethdbmax", "Database max-open.").
		Default("100").
		IntVar(&conf.EthComMysqlMaxOpen)
	app.Flag("reqlog", "Request log, support file only").
		PlaceHolder("/tmp/wallet/eth-req.log").
		StringVar(&conf.CmReqlog)
	app.Flag("debug", "Enable debug mode").Default("false").
		BoolVar(&conf.CmDebug)
	app.Flag("gethendpoint", "Geth endpoint").PlaceHolder("http://IP:PORT").Required().StringVar(&conf.EthComGethEndpoint)
	app.Flag("disable-auto-erc20", "Disable auto erc20 token flag").Default("false").BoolVar(&conf.DisableAutoErc20)
	conf.Nsqconfig = &nsqs.SimpleConfig{}
	app.Flag("nsqd", "Nsqd node address").Default("127.0.0.1:4151").
		StringVar(&conf.Nsqconfig.NsqAddress)
	app.Flag("nsqmif", "Nsqd max-in-flight").Default("100").
		IntVar(&conf.Nsqconfig.MaxInFlight)
	app.Flag("nsqlookups", "Nsq lookups address").Default("127.0.0.1:4161").
		StringsVar(&conf.Nsqconfig.Lookups)
	return conf
}
