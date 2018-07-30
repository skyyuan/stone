package chaos

import (
	"strconv"
	"strings"

	"stone/common"
	"stone/common/ethcomm"
	"stone/service/eth"
	"stone/service/web"
)

// GetOrCreateErc20CoinType get or create Erc20CoinType
func GetOrCreateErc20CoinType(webdb *common.GormDB, contractAddress string) (*web.Erc20CoinType, error) {
	erc20Coin := &web.Erc20CoinType{}
	var err error
	dbObj := webdb.Where("contract_address = ?", contractAddress).First(erc20Coin)
	if dbObj.Error != nil && !dbObj.RecordNotFound() {
		return nil, dbObj.Error
	}
	if dbObj.RecordNotFound() {
		common.Logger.Info("create erc20 address: ", contractAddress)
		erc20Coin, err = eth.GetERC20Info(contractAddress)
		if err != nil {
			return nil, err
		}
		erc20Coin.Symbol, err = getTokenSymbolName(webdb, erc20Coin.OriginalSymbol)
		if err != nil {
			return nil, err
		}
		err = webdb.Create(erc20Coin).Error
		if err != nil {
			return nil, err
		}
		err = webdb.Create(&web.EthereumEstimateGas{ContractAddress: contractAddress, EstimateGas: eth.DefaultErc20GasPrice}).Error
		if err != nil {
			return nil, err
		}
	}
	return erc20Coin, nil
}

// CreateUserErc20Token Create User Erc20Token
func CreateUserErc20Token(webdb *common.GormDB, address string, contractAddress string) error {
	userAddress := web.UserAddress{}
	dbObj := webdb.Where("address = ?", address).First(&userAddress)
	if dbObj.Error != nil && !dbObj.RecordNotFound() {
		return dbObj.Error
	}
	if dbObj.RecordNotFound() {
		return nil
	}
	userErc20Coin := &web.UserErc20Coin{Address: address, ContractAddress: contractAddress}
	return webdb.Where(*userErc20Coin).FirstOrCreate(userErc20Coin).Error
}

// SyncUserErc20Token Sync UserErc20Token
func SyncUserErc20Token(ethdb *ethcomm.GormDB, webdb *common.GormDB, blockNumber int64) error {
	var erc20ContractAddresses []string
	err := ethdb.Table("transaction_infos").Select(`distinct(erc20_contract_address)`).Where("block_number = ? and erc20_contract_address != \"\"", blockNumber).Pluck("erc20_contract_address", &erc20ContractAddresses).Error
	if err != nil {
		return err
	}
	for _, contractAddr := range erc20ContractAddresses {
		_, err := GetOrCreateErc20CoinType(webdb, contractAddr)
		if err != nil {
			return err
		}
	}

	var transactions []eth.TransactionInfo
	err = ethdb.Where("block_number = ?", blockNumber).Where(`erc20_contract_address != ""`).Find(&transactions).Error
	if err != nil {
		return err
	}
	transactions = removeDuplicateTransactions(transactions)
	for _, item := range transactions {
		err = CreateUserErc20Token(webdb, item.To, item.ERC20ContractAddress)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeDuplicateTransactions(txnInfos []eth.TransactionInfo) []eth.TransactionInfo {
	res := []eth.TransactionInfo{}
	encountered := map[string]string{}
	for _, item := range txnInfos {
		key := item.To + item.ERC20ContractAddress
		if _, ok := encountered[key]; !ok {
			encountered[key] = key
			res = append(res, item)
		}
	}
	return res
}

func getTokenSymbolName(webdb *common.GormDB, symbol string) (string, error) {
	coinType := web.Erc20CoinType{}
	dbObj := webdb.Where(`original_symbol  = ?`, symbol).Last(&coinType)
	if dbObj.RecordNotFound() {
		return symbol, nil
	}
	if dbObj.Error != nil && !dbObj.RecordNotFound() {
		return "", dbObj.Error
	}
	if coinType.OriginalSymbol == coinType.Symbol {
		return symbol + "_1", nil
	}

	strArr := strings.Split(coinType.Symbol, "_")
	strNum, err := strconv.Atoi(strArr[len(strArr)-1])
	if err != nil {
		return "", err
	}
	return symbol + "_" + strconv.Itoa(strNum+1), nil
}
