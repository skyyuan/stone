package eth

import (
	"fmt"
	"strconv"
	"strings"

	"stone/common"
	"stone/service/utils"
)

// GetGasPrice get gas price
func GetGasPrice() (*GasPrice, error) {
	var resp string
	err := endpointsManager.RPC(&resp, "eth_gasPrice")
	if err != nil {
		common.Logger.Debug(err)
		return nil, err
	}
	i256, ok := utils.ParseBig256(resp)
	if !ok {
		return nil, fmt.Errorf("Invalid hexstr %s", resp)
	}
	return &GasPrice{GasPrice: i256.Text(10)}, err
}

// GetTransactionCount get a user address transaction count
func GetTransactionCount(userAddress string) (string, error) {
	var resp string
	err := endpointsManager.RPC(&resp, "eth_getTransactionCount", userAddress, "latest")
	if err != nil {
		common.Logger.Debug(err)
		return HexPrefix, err
	}

	txCount, err := utils.HexToInt(resp)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(txCount, 10), err
}

// BlockNumber recent block number
func BlockNumber() (string, error) {
	var resp string
	err := endpointsManager.RPC(&resp, "eth_blockNumber")
	if err != nil {
		common.Logger.Debug(err)
		return HexPrefix, err
	}
	return resp, err
}
func DetailNodeInfo() ([]*EthereumNodeInfo, error) {
	var resp []*EthereumNodeInfo
	resp = endpointsManager.GetEndPoints()
	return resp, nil
}

// CliqueSnapshot recent block number
func CliqueSnapshot(blockNum string) (map[string]interface{}, error) {
	snapshot := make(map[string]interface{})

	err := endpointsManager.RPC(&snapshot, "clique_getSnapshot", blockNum)
	if err != nil {
		common.Logger.Debug(err)
		return snapshot, err
	}
	return snapshot, err
}

// GetCode if address is contract address return binary code else return 0x
func GetCode(address string) (string, error) {
	var resp string
	err := endpointsManager.RPC(&resp, "eth_getCode", address, "latest")
	if err != nil {
		common.Logger.Debug(err)
		return HexPrefix, err
	}
	return resp, err
}

// TxnStatus TxnStatus
func TxnStatus(txninfo *TransactionInfo, receipt *TransactionReceiptInfo) (int, error) {
	if receipt.BlockNumber == "" {
		return StatusPending, nil
	}

	if receipt.GasUsed == "" {
		return StatusPending, nil
	}

	gasUsed, err := utils.HexToInt(receipt.GasUsed)
	if err != nil {
		return StatusFailed, err
	}
	gas, err := strconv.ParseInt(txninfo.Gas, 10, 64)
	if err != nil {
		return StatusFailed, err
	}
	if gasUsed < gas || (gasUsed == EthGasLimit && gas == EthGasLimit) {
		return StatusPass, nil
	}
	return StatusFailed, nil
}

// GetComfirmedNum 根据transaction获取区块确认数
func GetComfirmedNum(txninfo *TransactionInfo) (int64, error) {
	if txninfo.Status == StatusPending {
		return 0, nil
	}
	currentBlockNumStr, err := BlockNumber()
	if err != nil {
		return 0, err
	}
	currentBlockNum, err := utils.HexToInt(currentBlockNumStr)
	if err != nil {
		return 0, err
	}
	delta := currentBlockNum - int64(txninfo.BlockNumber) //blockNum
	if delta < ConfirmedNum {
		return delta, nil
	}
	return ConfirmedNum, nil
}

// IsContract IsContract
func IsContract(address string) (bool, error) {
	code, err := GetCode(address)
	if err != nil {
		return false, err
	}
	// code is `0x`
	if len(code) == 2 {
		return false, nil
	}
	return true, nil
}

// IsERC20Txn is erc20 txn
func IsERC20Txn(txninfo *TransactionInfo) bool {
	return strings.HasPrefix(txninfo.Input, HexPrefix+ERC20MethodTransfer)
}

func truncationZero(str string) string {
	begin := 0
	for idx, item := range str {
		if item != '0' {
			begin = idx
			break
		}
	}
	return str[begin:]

}

// GenERC20Txn generate erc20 txn
func GenERC20Txn(txn *TransactionInfo) {
	txn.ERC20ContractAddress = txn.To
	if len(txn.Input) < ERC20TransferLength {
		common.Logger.Debug(txn.From, txn.To, txn.Input)
		return
	}
	txn.To = HexPrefix + txn.Input[34:74]

	i256, ok := utils.ParseBig256(HexPrefix + truncationZero(txn.Input[74:]))
	if ok {
		txn.Value = i256.Text(10)
	}
}
