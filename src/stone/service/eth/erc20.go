package eth

import (
	"encoding/hex"
	"strconv"
	"strings"

	"stone/common"
	"stone/service/web"
)

// GetERC20Info GetERC20Info
func GetERC20Info(contractAddress string) (*web.Erc20CoinType, error) {
	coinType := &web.Erc20CoinType{}
	// symbol
	var symbol string
	params := map[string]string{"to": contractAddress, "data": ERC20Symbol}
	err := endpointsManager.RPC(&symbol, "eth_call", params, "latest")
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}
	coinType.OriginalSymbol, err = convertToASCII(symbol)
	if err != nil {
		common.Logger.Error("cover symbol error: ", contractAddress, " err:", err, " abi:", symbol)
		return nil, err
	}
	// decimal
	params["data"] = ERC20Decimals
	var decimals string
	err = endpointsManager.RPC(&decimals, "eth_call", params, "latest")
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}
	coinType.Decimals, err = convertToDecimal(decimals)
	if err != nil {
		return nil, err
	}

	// name
	params["data"] = ERC20Name
	var name string
	err = endpointsManager.RPC(&name, "eth_call", params, "latest")
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}
	coinType.Name, err = convertToASCII(name)
	if err != nil {
		common.Logger.Error("cover name error: ", contractAddress, " err:", err, " abi:", name)
		return nil, err
	}

	coinType.ContractAddress = contractAddress
	coinType.Icon = DefaultErc20Icon
	return coinType, nil
}

func convertToDecimal(abiString string) (uint8, error) {
	tmpStr := strings.TrimLeft(abiString[2:], "0")
	if len(tmpStr) <= 0 {
		return uint8(0), nil
	}
	decimalsInt, err := strconv.ParseInt(tmpStr, 16, 8)
	if err != nil {
		return uint8(0), err
	}
	return uint8(decimalsInt), nil
}

func convertToASCII(abiString string) (string, error) {
	if len(abiString) < ERC20AbiDefaultLength {
		return "", nil
	}
	trimedAbi := strings.TrimRight(abiString[66+64:], "0")
	if len(trimedAbi)%2 != 0 {
		trimedAbi += "0"
	}
	res, err := hex.DecodeString(trimedAbi)
	if err != nil {
		common.Logger.Error(err)
		return "", err
	}
	return string(res), nil
}
