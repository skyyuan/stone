package v2

import (
	"github.com/labstack/echo"
	"stone/common"
	"stone/service/eth"
)

// RegisterAPI 注册v2版本的API
func RegisterAPI(router *echo.Echo) {
	v2 := router.Group("/v2")

	v2.POST("/get_txdetails", getTxnDetails)

}

func getTxnDetails(c echo.Context) error {
	params := new(EthereumTransactionID)
	if err := c.Bind(params); err != nil {
		return err
	}

	txDetails, err := eth.DetailTransactionInfoFromDB(params.TxID)
	if err != nil {
		return err
	}
	return common.JSONReturns(c, txDetails)
}
