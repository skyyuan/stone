package web

import (
	"strings"

	"stone/common"
)

// GetErc20CoinTypeList returns a list of erc20 coin types supported
func GetErc20CoinTypeList() []Erc20CoinType {
	var erc20s []Erc20CoinType

	db := common.DBBegin()
	defer db.DBRollback()

	db.Where(`is_visible = 1`).Order("weight desc, id asc").Find(&erc20s)

	return erc20s
}

// GetEstimateGas return the estimate Gas for a call or transaction.
// on given contract address
//
// contractAddress is "" when ETH.
func GetEstimateGas(contractAddress string) *EthereumEstimateGas {
	var estimateGas EthereumEstimateGas

	db := common.DBBegin()
	defer db.DBRollback()

	db.Where("contract_address = ?", strings.ToLower(contractAddress)).First(&estimateGas)

	return &estimateGas
}

// GetAllEstimateGases return all estimate Gases for a call or transaction
// on all contract address
func GetAllEstimateGases() (estimateGases []EthereumEstimateGas) {
	db := common.DBBegin()
	defer db.DBRollback()

	db.Find(&estimateGases)

	return
}
