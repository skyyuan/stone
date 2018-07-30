package chaos

import (
	"testing"

	"stone/common"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetSymbolName(t *testing.T) {
	mocket.Catcher.Register()
	mockDb, err := gorm.Open(mocket.DRIVER_NAME, "mock DB")
	mockDb.LogMode(true)
	if err == nil {
		common.DB = mockDb
	}

	t.Run("eth tx and abnormal testcase", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  "SELECT count(*) FROM \"erc20_coin_types\"  WHERE \"erc20_coin_types\".\"deleted_at\" IS NULL",
				Response: []map[string]interface{}{{"total": 1}},
			},
			{

				Pattern: "SELECT * FROM \"erc20_coin_types\"  WHERE \"erc20_coin_types\".\"deleted_at\" IS NULL AND ((original_symbol  = ATM)) ORDER BY \"erc20_coin_types\".\"id\" DESC LIMIT 1",
				Response: []map[string]interface{}{
					{
						"symbol":           "ATM_1",
						"original_symbol":  "ATM",
						"contract_address": "0xad5c52d34b5e4cf0c030ff5726ba44b2c4802e22",
						"name":             "atm",
						"decimals":         1,
						"weight":           1,
						"is_visible":       false,
						"is_suicided":      false,
						"icon":             "ddd",
					},
				},
			},
		})
	})

	db := common.DBBegin()
	defer db.DBRollback()
	res, _ := getTokenSymbolName(db, "ATM")
	assert.Equal(t, res, "ATM_2")
}
