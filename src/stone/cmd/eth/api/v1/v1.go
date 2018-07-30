package v1

import (
	"fmt"
	"strconv"

	ecom "github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo"
	"stone/common"
	"stone/service/eth"
	"stone/service/utils"
)

func txnlist(c echo.Context) error {
	p := EthereumTransactionListPayload{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if err := c.Echo().Validator.Validate(p); err != nil {
		return err
	}
	/* no more check
	if !ecom.IsHexAddress(p.UserAddress) ||
		(p.ContractAddress != "" && !ecom.IsHexAddress(p.ContractAddress)) {
		return common.BizError1001
	}
	*/
	if p.PerPage == 0 {
		p.PerPage = 20
	}
	txlist, total := eth.EthereumTransactionList(p.UserAddress, p.ContractAddress,
		p.Page, p.PerPage)

	return common.JSONReturns(c, txlist, p.Page, total, p.PerPage)
}

func txnlistByBlock(c echo.Context) error {
	p := EthereumBlockPayload{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if err := c.Echo().Validator.Validate(p); err != nil {
		return err
	}

	if p.PerPage == 0 {
		p.PerPage = 20
	}
	txlist, total := eth.EthereumTransactionListByBlock(p.BlockNumber,
		p.Page, p.PerPage)

	return common.JSONReturns(c, txlist, p.Page, total, p.PerPage)
}
func getBalance(c echo.Context) error {
	params := new(EthereumAccount)
	if err := c.Bind(params); err != nil {
		return err
	}

	if !ecom.IsHexAddress(params.UserAddress) {
		return common.BizError1001
	}

	var balance1 *eth.BalanceAmount
	var balance2 *eth.BalanceAmount
	mbalance := make(map[string]string)
	var err error
	balance1, err = eth.EthereumBalance(params.UserAddress)
	if err != nil {
		return err
	}
	mbalance["ATM"] = balance1.Amount
	if params.ContractAddress != "" && ecom.IsHexAddress(params.ContractAddress) {
		balance2, err = eth.Erc20Balance(params.UserAddress, params.ContractAddress)
		if err != nil {
			return err
		}
		mbalance[params.ContractAddress] = balance2.Amount
	}

	fmt.Println("params.ContractAddress: ", params.ContractAddress, "params.UserAddress:", params.UserAddress, "balance:", balance1.Amount)
	return common.JSONReturns(c, mbalance)
}

func getMultipleBalances(c echo.Context) error {
	params := new(EthereumAccountMultipleContract)
	if err := c.Bind(params); err != nil {
		return err
	}

	if !ecom.IsHexAddress(params.UserAddress) {
		return common.BizError1001
	}

	mbalance := make(map[string]string)
	for _, contract := range params.ContractAddresses {
		var err error
		var balance *eth.BalanceAmount
		if contract == "" {
			balance, err = eth.EthereumBalance(params.UserAddress)
		} else if ecom.IsHexAddress(contract) {
			balance, err = eth.Erc20Balance(params.UserAddress, contract)
		} else {
			return common.BizError1001
		}
		if err != nil {
			return err
		}

		mbalance[contract] = balance.Amount
	}

	return common.JSONReturns(c, mbalance)
}

func getTxnDetails(c echo.Context) error {
	params := new(EthereumTransactionHash)
	if err := c.Bind(params); err != nil {
		return err
	}

	var txDetails *eth.TransactionInfo
	var err error
	if len(params.TxHash) == 2+2*ecom.HashLength && ecom.IsHex(params.TxHash) {
		txDetails, err = eth.DetailTransactionInfo(params.TxHash)
		if err != nil {
			return err
		}
		return common.JSONReturns(c, txDetails)
	}
	return common.BizError1001
}

func getGethNodeInfo(c echo.Context) error {
	var nodeDetails []*eth.EthereumNodeInfo
	var err error
	nodeDetails, err = eth.DetailNodeInfo()
	fmt.Println("nodeDetails, err = ", nodeDetails, err)
	if err != nil {
		return err
	}
	return common.JSONReturns(c, nodeDetails)

}
func getCliqueSnapshot(c echo.Context) error {
	lastest, err := eth.BlockNumber()
	if err != nil {
		common.Logger.Debug(err)
		return err
	}
	snapshot, err := eth.CliqueSnapshot(lastest)
	if err != nil {
		return err
	}
	return common.JSONReturns(c, snapshot)
}

func getGasPrice(c echo.Context) error {
	gasPrice, err := eth.GetGasPrice()
	if err != nil {
		return err
	}

	return common.JSONReturns(c, gasPrice)
}

func sendRawTransaction(c echo.Context) error {
	var params EthereumSendRawTransactionPayload
	if err := c.Bind(&params); err != nil {
		return err
	}

	var rsp EthereumSendRawTransaction
	if len(params.SignedTransaction) > 2 && params.SignedTransaction[:2] == "0x" {
		r, err := eth.SendRawTransaction(params.SignedTransaction)
		if err != nil {
			return err
		}
		rsp.Hash = r
		return common.JSONReturns(c, rsp)
	}
	return common.BizError1001
}

func getTxnCount(c echo.Context) error {
	var params EthereumTxnCountPayload
	if err := c.Bind(&params); err != nil {
		return err
	}
	resp, err := eth.GetTransactionCount(params.UserAddress)
	if err != nil {
		return err
	}
	return common.JSONReturns(c, resp)
}

func getTxnReceipt(c echo.Context) error {
	var params EthereumTransactionHash
	if err := c.Bind(&params); err != nil {
		return err
	}
	if len(params.TxHash) == 2+2*ecom.HashLength && ecom.IsHex(params.TxHash) {
		resp, err := eth.TransactionReceipt(params.TxHash)
		if err != nil {
			return err
		}
		return common.JSONReturns(c, resp)
	}
	return common.BizError1001
}
func getBlockByBlockNumber(c echo.Context) error {
	var params EthereumBlockPayload
	if err := c.Bind(&params); err != nil {
		return err
	}
	lastest, err := eth.BlockNumber()
	if err != nil {
		common.Logger.Debug(err)
		return err
	}
	v1, _ := strconv.Atoi(params.BlockNumber)
	v2, _ := utils.HexToInt(lastest)
	if v1 > int(v2) {
		fmt.Println("params.BlockNumber > lastest ", params.BlockNumber, ",", lastest)
		return common.BizError1001
	}
	blockInfo := eth.EthereumBlock(params.BlockNumber)
	return common.JSONReturns(c, blockInfo)
}
func blocklist(c echo.Context) error {
	p := EthereumBlockListPayload{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if err := c.Echo().Validator.Validate(p); err != nil {
		return err
	}

	if p.PerPage == 0 {
		p.PerPage = 20
	}
	blocklist, total := eth.EthereumBlockList(p.Page, p.PerPage)

	return common.JSONReturns(c, blocklist, p.Page, total, p.PerPage)
}
func getATMPrice(c echo.Context) error {
	p := EthereumATMpricePayload{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	price, err := eth.HttpGetATMPrice()
	if err != nil {
		return err
	}

	return common.JSONReturns(c, price)
}
