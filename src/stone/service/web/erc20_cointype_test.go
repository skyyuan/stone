package web

import (
	"testing"

	"stone/common"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
)

func TestGetErc20CoinTypeList(t *testing.T) {
	mocket.Catcher.Register()
	db, err := gorm.Open(mocket.DRIVER_NAME, "mocket database erc20_coin_types")
	if err != nil {
		t.Errorf("Failed to open database. %v\n", err)
	}
	common.DB = db

	t.Run("erc20 coin types", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: "SELECT * FROM \"erc20_coin_types\"  WHERE \"erc20_coin_types\".\"deleted_at\" IS NULL",
				Response: []map[string]interface{}{
					{
						"name":             "Digix DAO",
						"symbol":           "DGD",
						"icon_url":         "https://files.coinmarketcap.com/static/img/coins/32x32/digixdao.png",
						"contract_address": "0x95f95c3d758a3c488da74d572513dfe5eec1f930",
						"decimals":         8,
					},
					{
						"name":             "Ethereum",
						"symbol":           "ETH",
						"icon_url":         "http://www.baidu.com",
						"contract_address": "",
						"decimals":         2,
					},
				},
			},
		})

		coins := GetErc20CoinTypeList()
		if len(coins) != 2 {
			t.Errorf("count of coin types does not match")
			return
		}
		if "Digix DAO" != coins[0].Name || "DGD" != coins[0].Symbol ||
			"0x95f95c3d758a3c488da74d572513dfe5eec1f930" != coins[0].ContractAddress ||
			8 != coins[0].Decimals {
			t.Errorf("unexpected coin type")
		}
		if "Ethereum" != coins[1].Name || "ETH" != coins[1].Symbol ||
			"" != coins[1].ContractAddress || 2 != coins[1].Decimals {
			t.Errorf("unexpected coin type")
		}

	})
}

func TestGetEstimateGas(t *testing.T) {
	mocket.Catcher.Register()
	db, err := gorm.Open(mocket.DRIVER_NAME, "mocket database ethereum_estimate_gas")
	if err != nil {
		t.Errorf("Failed to open database. %v\n", err)
	}
	common.DB = db

	t.Run("ethereum estimate gas", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: "SELECT * FROM \"ethereum_estimate_gas\"  WHERE \"ethereum_estimate_gas\".\"deleted_at\" IS NULL",
				Response: []map[string]interface{}{
					{
						"contract_address": "",
						"estimate_gas":     "21000",
					},
				},
			},
		})

		gas := GetEstimateGas("")
		if "" != gas.ContractAddress || 21000 != gas.EstimateGas {
			t.Errorf("unexpected estimate gas")
		}
	})

	t.Run("erc20 estimate gas", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: "SELECT * FROM \"ethereum_estimate_gas\"  WHERE \"ethereum_estimate_gas\".\"deleted_at\" IS NULL AND ((contract_address = 0x5f81dc51bdc05f4341afbfa318af5d82c607acad)) ORDER BY \"ethereum_estimate_gas\".\"id\" ASC LIMIT 1",
				Response: []map[string]interface{}{
					{
						"contract_address": "0x5f81dc51bdc05f4341afbfa318af5d82c607acad",
						"estimate_gas":     "53000",
					},
				},
			},
		})

		gas := GetEstimateGas("0x5f81DC51BDC05f4341afbfa318af5d82c607acad")
		if "0x5f81dc51bdc05f4341afbfa318af5d82c607acad" != gas.ContractAddress ||
			53000 != gas.EstimateGas {
			t.Errorf("unexpected estimate gas")
		}
	})

	t.Run("non-existing coin type", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: "SELECT * FROM \"ethereum_estimate_gas\"  WHERE \"ethereum_estimate_gas\".\"deleted_at\" IS NULL AND ((contract_address = )) ORDER BY \"ethereum_estimate_gas\".\"id\" ASC LIMIT 1",
				Response: []map[string]interface{}{
					{
						"contract_address": "",
						"estimate_gas":     "21000",
					},
				},
			},
		})

		gas := GetEstimateGas("0xa848a7FECfd327f1e1595EcD1A8A5E6FF3FC9113")
		if "" != gas.ContractAddress || 0 != gas.EstimateGas {
			t.Errorf("unexpected estimate gas")
		}
	})
}

func TestGetAllEstimateGases(t *testing.T) {
	mocket.Catcher.Register()
	db, err := gorm.Open(mocket.DRIVER_NAME, "mocket database ethereum_estimate_gas")
	if err != nil {
		t.Errorf("Failed to open database. %v\n", err)
	}
	common.DB = db

	t.Run("all estimate gas", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern: "SELECT * FROM \"ethereum_estimate_gas\"  WHERE \"ethereum_estimate_gas\".\"deleted_at\" IS NULL",
				Response: []map[string]interface{}{
					{
						"contract_address": "",
						"estimate_gas":     "21000",
					},
					{
						"contract_address": "0x5f81dc51bdc05f4341afbfa318af5d82c607acad",
						"estimate_gas":     "53000",
					},
				},
			},
		})

		gases := GetAllEstimateGases()
		if len(gases) != 2 {
			t.Errorf("count of coins is wrong")
			return
		}
		if "" != gases[0].ContractAddress || 21000 != gases[0].EstimateGas ||
			"0x5f81dc51bdc05f4341afbfa318af5d82c607acad" != gases[1].ContractAddress ||
			53000 != gases[1].EstimateGas {
			t.Errorf("unexpected estimate gas")
		}
	})
}
