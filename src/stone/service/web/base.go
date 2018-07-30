package web

import "github.com/jinzhu/gorm"

// Erc20CoinType represents an erc20 coin type
type Erc20CoinType struct {
	gorm.Model
	Name            string `gorm:"type:varchar(255);not null" json:"name"`
	Symbol          string `gorm:"type:varchar(20);index:symbol" json:"symbol"`
	Icon            string `gorm:"type:varchar(255)" json:"iconUrl"`
	ContractAddress string `gorm:"type:varchar(42);unique;not null" json:"contractAddress"`
	Decimals        uint8  `gorm:"type:tinyint unsigned;not null;default:0" json:"decimals"`
	Weight          uint   `gorm:"type:int unsigned;not null;default:0;index" json:"-"`
	IsVisible       bool   `gorm:"not null;default:false;index" json:"-"`
	IsSuicided      bool   `gorm:"not null;default:false;index" json:"-"`
	OriginalSymbol  string `gorm:"type:varchar(20);index:original_symbol" json:"-"`
}

// EthereumEstimateGas represents estimate Gas used when make a call or tx
type EthereumEstimateGas struct {
	gorm.Model
	// ContractAddress is "" when ETH
	ContractAddress string `gorm:"type:varchar(42);unique;not null"`
	EstimateGas     int64  `gorm:"type:bigint unsigned;not null;default:0"`
}

// Device device info
type Device struct {
	gorm.Model
	DeviceType uint   `gorm:"type:int unsigned;not null;default:0;unique_index:idx_device_type_uuid" json:"device_type"`
	UUID       string `gorm:"type:varchar(100);unique_index:idx_device_type_uuid"`
}

// UserAddress contains user registered address
type UserAddress struct {
	gorm.Model
	DeviceID    uint   `gorm:"type:int unsigned;not null;unique_index:idx_address_type_address_deviceid" json:"device_id"`
	Address     string `gorm:"type:varchar(200);index:idx_address_type_address;unique_index:idx_address_type_address_deviceid" json:"address"`
	AddressType uint   `gorm:"type:int unsigned;not null;default:0;index:idx_address_type_address;unique_index:idx_address_type_address_deviceid" json:"address_type"`
}

// UserErc20Coin UserErc20Coin
type UserErc20Coin struct {
	gorm.Model
	Address         string `gorm:"type:varchar(42);unique_index:idx_address_contract_address"`
	ContractAddress string `gorm:"type:varchar(42);unique_index:idx_address_contract_address"`
	IsNew           bool   `gorm:"not null;default:true;index" json:"-"`
}
