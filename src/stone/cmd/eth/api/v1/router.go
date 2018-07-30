package v1

import (
	"github.com/labstack/echo"
)

// RegisterAPI 注册v1版本的API
func RegisterAPI(router *echo.Echo) {
	v1 := router.Group("/v1")

	v1.POST("/get_txhistory", txnlist)
	v1.POST("/get_txhistoryByBlock", txnlistByBlock)
	v1.POST("/get_balance", getBalance)
	v1.POST("/get_txdetails", getTxnDetails)
	v1.POST("/get_multiplebalances", getMultipleBalances)
	v1.POST("/get_gasprice", getGasPrice)
	v1.POST("/send_raw_transaction", sendRawTransaction)
	v1.POST("/get_txcount", getTxnCount)
	v1.POST("/get_txn_receipt", getTxnReceipt)
	v1.POST("/get_blocklist", blocklist)
	v1.POST("/get_block", getBlockByBlockNumber)
	v1.POST("/get_ATMprice", getATMPrice)
	v1.POST("/get_CliqueSnapshot", getCliqueSnapshot)
	v1.POST("/get_GethNodeInfo", getGethNodeInfo)	
}
