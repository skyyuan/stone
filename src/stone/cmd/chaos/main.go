package main

import (
	"os"
	"os/signal"
	"syscall"

	"stone/cmd/chaos/runner"

	"github.com/alecthomas/kingpin"
	"github.com/labstack/gommon/log"
	"stone/common"
	"stone/common/chaoscomm"
	"stone/common/ethcomm"
	"stone/nsqs"
	"stone/service/eth"
)

var (
	// Version version
	Version = "0.0.1"
	app     = kingpin.New("app", "Chaos").DefaultEnvars()
	cmdRun  = app.Command("run", "Run application").Default()
)

func main() {
	conf := configure()
	kingpin.Version(Version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	eth.ServiceInit(conf)
	ethcomm.InitDB(conf)
	defer ethcomm.DBClose()
	common.InitDB(conf)
	defer common.DBClose()
	common.Logger = log.New("Chaos")
	if conf.Debug() {
		common.Logger.SetLevel(log.DEBUG)
	}

	err := nsqs.InitConfig(conf.Nsqconfig)
	if err != nil {
		common.Logger.Fatal(err)
	}

	runner.Register()

	nsqs.Start()

	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal

	eth.ServiceDone()
	nsqs.Stop()
}

func configure() *chaoscomm.ChaosConfig {
	conf := &chaoscomm.ChaosConfig{}
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
	app.Flag("ethdb", "Ethereum database connection URL, only support mysql.").
		PlaceHolder("USER:PWD@tcp(DBURL:DBPORT)/DBSCHEMA?charset=utf8&parseTime=True&loc=Local").
		Required().
		StringVar(&conf.EthComMysqlURL)
	app.Flag("ethdbidle", "Ethereum database idel connection numbers.").
		Default("10").
		IntVar(&conf.EthComMysqlIdle)
	app.Flag("ethdbmax", "Ethereum database max-open.").
		Default("100").
		IntVar(&conf.EthComMysqlMaxOpen)
	app.Flag("reqlog", "Request log, support file only").
		PlaceHolder("/tmp/wallet/eth-req.log").
		StringVar(&conf.CmReqlog)
	app.Flag("debug", "Enable debug mode").Default("false").
		BoolVar(&conf.CmDebug)
	app.Flag("gethendpoint", "Geth endpoint").PlaceHolder("http://IP:PORT").Required().StringVar(&conf.EthComGethEndpoint)
	conf.Nsqconfig = &nsqs.SimpleConfig{}
	app.Flag("nsqd", "Nsqd node address").Default("127.0.0.1:4151").
		StringVar(&conf.Nsqconfig.NsqAddress)
	app.Flag("nsqmif", "Nsqd max-in-flight").Default("100").
		IntVar(&conf.Nsqconfig.MaxInFlight)
	app.Flag("nsqlookups", "Nsq lookups address").Default("127.0.0.1:4161").
		StringsVar(&conf.Nsqconfig.Lookups)
	return conf
}
