package web

import (
	"stone/common"
	"stone/nsqs"
)

// RegisterAddress register Address
func RegisterAddress(address string, addressType uint, deviceType uint, deviceUUID string) error {
	db := common.DBBegin()
	defer db.DBRollback()

	device := &Device{DeviceType: deviceType, UUID: deviceUUID}
	err := db.Where(*device).FirstOrCreate(device).Error
	if err != nil {
		return err
	}

	var userAddress UserAddress
	dbObj := db.Where("address = ? and address_type = ?", address, addressType).First(&userAddress)
	if dbObj.Error != nil && !dbObj.RecordNotFound() {
		return dbObj.Error
	}
	if dbObj.RecordNotFound() {
		err := nsqs.PostTopic("eth_erc20_from_new_user", address)
		if err != nil {
			return err
		}
	}

	userAddress = UserAddress{
		DeviceID:    device.ID,
		Address:     address,
		AddressType: addressType}
	err = db.Where(userAddress).FirstOrCreate(&userAddress).Error
	if err != nil {
		return err
	}
	db.DBCommit()
	return nil
}

// UserErc20TokenList UserErc20Token List
func UserErc20TokenList(address string) ([]Erc20CoinType, error) {
	db := common.DBBegin()
	defer db.DBRollback()

	var contractAddresses []string
	var erc20Coins []Erc20CoinType
	err := db.Table("user_erc20_coins").Select(`contract_address`).Where(`address = ?`, address).Pluck("contract_address", &contractAddresses).Error
	if err != nil {
		return nil, err
	}
	err = db.Where(`contract_address in (?)`, contractAddresses).Find(&erc20Coins).Error
	if err != nil {
		return nil, err
	}
	var erc20s []Erc20CoinType
	err = db.Where(`is_visible = 1`).Order("weight desc, id asc").Find(&erc20s).Error
	if err != nil {
		return nil, err
	}
	erc20s = append(erc20s, erc20Coins...)
	return erc20s, nil
}

// FetchNewUserErc20TokenList User NewErc20 TokenList
func FetchNewUserErc20TokenList(address string) ([]Erc20CoinType, error) {
	db := common.DBBegin()
	defer db.DBRollback()
	var contractAddresses []string
	var erc20Coins []Erc20CoinType
	err := db.Table("user_erc20_coins").Select(`contract_address`).Where(`address = ? and is_new = true`, address).Pluck("contract_address", &contractAddresses).Error
	if err != nil {
		return nil, err
	}
	err = db.Where(`contract_address in (?)`, contractAddresses).Find(&erc20Coins).Error
	if err != nil {
		return nil, err
	}
	err = db.Table("user_erc20_coins").Where(`address = ?`, address).Where(`contract_address in (?)`, contractAddresses).Update("is_new", 0).Error
	db.DBCommit()
	return erc20Coins, err
}
