package main

import (
	"stone/common"
	"stone/common/auth"
	"stone/service/web"
)

func dDbMigrate2() {
	db := common.DBBegin()

	db.AutoMigrate(&web.Erc20CoinType{}, &web.EthereumEstimateGas{}, &auth.AppAuth{}, &web.UserAddress{}, &web.UserErc20Coin{}, &web.Device{})

	db.Commit()
}

func DbMigrate() {
	db := common.DBBegin()

	db.AutoMigrate(&web.Erc20CoinType{}, &web.EthereumEstimateGas{}, &auth.AppAuth{}, &web.UserAddress{}, &web.UserErc20Coin{}, &web.Device{})

	db.Commit()
}
