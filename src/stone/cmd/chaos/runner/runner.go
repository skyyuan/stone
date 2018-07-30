package runner

import (
	"encoding/json"
	"strconv"

	"github.com/nsqio/go-nsq"
	"stone/common"
	"stone/common/ethcomm"
	"stone/nsqs"
	"stone/service/chaos"
)

// Register register nsq listener
func Register() {
	nsqs.RegisterDefault("eth_erc20_from_new_user", syncERC20FromNewUser)
	nsqs.RegisterDefault("eth_sync_block_erc20", syncERC20FromBlock)
}

func syncERC20FromNewUser(m *nsq.Message) error {
	var address string
	err := json.Unmarshal(m.Body, &address)
	if err != nil {
		common.Logger.Error(err)
		return nil
	}
	db := ethcomm.DBBegin()
	defer db.DBRollback()

	var erc20Addresses []string
	err = db.Table("transaction_infos").Select(`DISTINCT(erc20_contract_address)`).Where("`to` = ? and erc20_contract_address != \"\"", address).Pluck("erc20_contract_address", &erc20Addresses).Error

	if err != nil {
		common.Logger.Info(err)
		return err
	}

	webdb := common.DBBegin()
	defer webdb.DBRollback()

	for _, erc20Address := range erc20Addresses {
		_, err := chaos.GetOrCreateErc20CoinType(webdb, erc20Address)
		if err != nil {
			return err
		}
		err = chaos.CreateUserErc20Token(webdb, address, erc20Address)
		if err != nil {
			return err
		}
	}
	webdb.DBCommit()
	return nil
}

func syncERC20FromBlock(m *nsq.Message) error {
	blockNum, err := strconv.ParseInt(string(m.Body), 10, 64)
	if err != nil {
		return err
	}

	ethdb := ethcomm.DBBegin()
	defer ethdb.DBRollback()
	webdb := common.DBBegin()
	defer webdb.DBRollback()

	err = chaos.SyncUserErc20Token(ethdb, webdb, blockNum)
	if err != nil {
		common.Logger.Error("syncERC20FromBlock: ", blockNum, " [", err, "]")
		return err
	}

	webdb.DBCommit()
	return nil
}
