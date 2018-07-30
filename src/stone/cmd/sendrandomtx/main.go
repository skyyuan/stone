package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/gommon/log"
	"stone/common"
	"stone/service/eth"
)

var endpointsManager1 = eth.NewEndPointsManager()
var endpointsManager2 = eth.NewEndPointsManager()
var endpointsManager3 = eth.NewEndPointsManager()

func GetAccounts(end *eth.EndpointsManager) ([]string, error) {
	var resp []string

	err := end.RPC(&resp, "eth_accounts")
	if err != nil {
		common.Logger.Debug(err)
		return []string{}, err
	}
	return resp, nil
}
func SendTransaction(end *eth.EndpointsManager, from, to, value string) (string, error) {
	var resp string
	params := map[string]string{"from": from, "to": to, "value": value, "gas": "0x76c0", "gasPrice": "0x9184e72a000", "data": ""}
	err := end.RPC(&resp, "eth_sendTransaction", params)
	if err != nil {
		common.Logger.Debug(err)
		return "", err
	}
	return resp, nil
}
func SendTx() {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if account, _ := GetAccounts(endpointsManager1); len(account) != 0 {
				SendTransaction(endpointsManager1, account[0], account[1], "0x10200000000000000")
				SendTransaction(endpointsManager1, account[0], account[1], "0x12000000000000000")
				SendTransaction(endpointsManager1, account[1], account[0], "0x5400000000000000")
				SendTransaction(endpointsManager1, account[1], account[0], "0x10200000000000000")
				SendTransaction(endpointsManager1, account[1], account[0], "0x1000000000000000")
				SendTransaction(endpointsManager1, account[0], account[1], "0x22000000000000000")
			}
			if account, _ := GetAccounts(endpointsManager2); len(account) != 0 {
				SendTransaction(endpointsManager2, account[0], account[1], "0x10200000000000000")
				SendTransaction(endpointsManager2, account[1], account[0], "0x10012000000000")
				SendTransaction(endpointsManager2, account[0], account[1], "0x1002220000000000")
				SendTransaction(endpointsManager2, account[1], account[0], "0x10200000000000000")
				SendTransaction(endpointsManager2, account[0], account[1], "0x10012000000000")
				SendTransaction(endpointsManager2, account[1], account[0], "0x1002220000000000")
			}
			if account, _ := GetAccounts(endpointsManager3); len(account) != 0 {
				SendTransaction(endpointsManager3, account[0], account[1], "0x1024000000000000")
				SendTransaction(endpointsManager3, account[0], account[1], "0x10400000000000000")
				SendTransaction(endpointsManager3, account[1], account[0], "0x1008000000000000")
				SendTransaction(endpointsManager3, account[1], account[0], "0x1024000000000000")
				SendTransaction(endpointsManager3, account[1], account[0], "0x1040000000000000")
				SendTransaction(endpointsManager3, account[0], account[1], "0x1008000000000000")
			}
		}
	}
}
func main() {
	common.Logger = log.New("send tick")
	endpointsManager1.AddEndPoint("http://47.254.147.68:8545", 1)
	endpointsManager2.AddEndPoint("http://47.88.228.248:8545", 1)
	endpointsManager3.AddEndPoint("http://47.89.245.117:8545", 1)

	go endpointsManager1.Run()
	go endpointsManager2.Run()
	go endpointsManager3.Run()

	go SendTx()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
