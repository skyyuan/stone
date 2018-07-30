package eth

import (
	"stone/common"
	"stone/service/utils"
)

// EthereumBalance get user ethereum balance
func EthereumBalance(userAddress string) (*BalanceAmount, error) {
	var resp string
	err := endpointsManager.RPC(&resp, "eth_getBalance", userAddress, "latest")
	if err != nil {
		common.Logger.Debug(err)
		return &BalanceAmount{}, err
	}

	amount := "0"
	t, ok := utils.ParseBig256(resp)
	if ok && t != nil {
		amount = t.Text(10)
	}
	return &BalanceAmount{Amount: amount}, err

}

// Erc20Balance get user balance in the contract token account
func Erc20Balance(userAddress string, contractAddress string) (*BalanceAmount, error) {
	var resp string
	params := map[string]string{"to": contractAddress, "data": utils.PaddingData(ERC20MethodBalanceOf, userAddress)}
	err := endpointsManager.RPC(&resp, "eth_call", params, "latest")
	if err != nil {
		common.Logger.Debug(err)
		return &BalanceAmount{}, err
	}

	amount := "0"
	t, ok := utils.ParseBig256(resp)
	if ok && t != nil {
		amount = t.Text(10)
	}
	return &BalanceAmount{Amount: amount}, err
}
