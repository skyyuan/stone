package eth

import (
	"encoding/json"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"stone/common/ethcomm"
	"stone/service/utils"
)

func init() {
	config := ethcomm.EthConfig{
		GethEndpoint: RopstenTestNode,
	}

	ServiceInit(&config)
}

func TestErc20DetailTransactionInfo(t *testing.T) {
	info, err := DetailTransactionInfo(
		"0xd2397b1e1a7b3230616e6232c2de205e1cf4d2870c7204f979a79f8ed11a9a90")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, info.From, "0x973cca58abebcad228a0cfa0bc91f8c3101090f0",
		"`from` address does not match")
	assert.Equal(t, info.IsContract, true,
		"IsContract does not match")
	assert.Equal(t, info.To, "0x54fb2b073926f20aa46604b00ccab89e50b5505f",
		"`To` address does not match")
	assert.Equal(t, info.BlockNumber, int64(1671225),
		"BlockNumber does not match")
	assert.Equal(t, info.BlockHash,
		"0xb418fff468d19a20fbf6337aa9c380626a8b412249b1157447571b01255e2c0c",
		"BlockHash does not match")
	assert.Equal(t, info.Gas, "36739",
		"Gas does not match")
	assert.Equal(t, info.GasPrice, "21000000000",
		"GasPrice does not match")
	assert.Equal(t, info.Value, "0",
		"Value does not match")
	assert.Equal(t, info.Status, 0,
		"Status does not match")
	assert.Equal(t, info.Input, "0xa9059cbb000000000000000000000000413b2b64fd25ab1b77d11"+
		"7005285dd4d40c509f70000000000000000000000000000000000000000000000000000000000000503",
		"Input Data does not match")

	// Invalid tx hash makes the service panic
}

func TestEthDetailTransactionInfo(t *testing.T) {
	info, err := DetailTransactionInfo(
		"0x4a64cb0154b97586620f3a8ae8e54655a37f858fae37f15418dbf06ef17e557c")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, info.From, "0x413b2b64fd25ab1b77d117005285dd4d40c509f7",
		"`from` address does not match")
	assert.Equal(t, info.IsContract, false,
		"IsContract does not match")
	assert.Equal(t, info.To, "0x973cca58abebcad228a0cfa0bc91f8c3101090f0",
		"`To` address does not match")
	assert.Equal(t, info.BlockNumber, int64(1671239),
		"BlockNumber does not match")
	assert.Equal(t, info.BlockHash,
		"0xc851ff4e05894b7e28a2145f2cd5c31292c0d8fb030454db61dbd2e23c666404",
		"BlockHash does not match")
	assert.Equal(t, info.Gas, "21000",
		"Gas does not match")
	assert.Equal(t, info.GasPrice, "21000000000",
		"GasPrice does not match")
	assert.Equal(t, info.Value, "2100000000000000",
		"Value does not match")
	assert.Equal(t, info.Status, 0,
		"Status does not match")
	assert.Equal(t, info.Input, "0x",
		"Input Data does not match")
	assert.Equal(t, info.Timestamp, int64(1505372757),
		"Status does not match")
	// Invalid tx hash makes the service panic
}

func TestTransactionReceipt(t *testing.T) {
	receipt, err := TransactionReceipt("0xcabe56275238a18e180f670570b0e81e2a9931a84b383804aa374a30ce538c4d")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, receipt.GasUsed, "0x7a120")
}

func TestGetBlockInfo(t *testing.T) {
	blockInfo, err := GetBlockInfo(utils.Int64ToHex(1658911))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, blockInfo.Number, int64(1658911), "blockNumber does not match")
	assert.Equal(t, blockInfo.Hash,
		"0x4ebc5e8a89c03de18323eb49b2a658fa714dcf3fa1d38dd5e52bfc2f3c289b9a",
		"blockNumber does not match")
	assert.Equal(t, blockInfo.ParentHash,
		"0xeb3060abf4d38f7365cdb2ff8fd0fae89d0a2d961071cb1a11d5b78dc60eea0d",
		"blockNumber does not match")
	assert.Equal(t, blockInfo.Timestamp, int64(1505200249), "timestamp does not match")

	txns := make([]map[string]interface{}, 0)
	err = json.Unmarshal(blockInfo.Transactions, &txns)
	txHashs := make([]string, 0, len(txns))
	for _, item := range txns {
		tmpTxn := &TransactionInfo{}
		MapAdapter(item, tmpTxn)
		txHashs = append(txHashs, tmpTxn.Hash)
	}
	assert.Equal(t, txHashs,
		[]string{"0x92b994181828a12189e2a3f09f4af0b9718b1092755688a1e5b777098d189fad",
			"0x3ec86434e6698be2751d4f67e783134fc3198c0a781c84bf290b2724c3060be2",
			"0x00d36d9a3cd1619f4bf80b530678f39bb1e6c297a9842acc3624441c7c7a7004",
			"0xfe8ed4ee236bc84a3b760db0c4364ae7eb0463581604ea92b8cfa51876608e24"},
		"tx list does not match")

	// blockNumber "0xffffffff" makes service panic
}

func TestEthereumTransactionList(t *testing.T) {
	mocket.Catcher.Register()
	db, err := gorm.Open(mocket.DRIVER_NAME, "mock DB")
	if err == nil {
		ethcomm.DB = db
	}

	t.Run("eth tx and abnormal testcase", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  "SELECT count(*) FROM \"transaction_infos\"  WHERE (erc20_contract_address = ) AND (`from` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155 OR `to` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155)",
				Response: []map[string]interface{}{{"total": 6}},
			},
			{
				Pattern: "SELECT * FROM \"transaction_infos\"  WHERE (erc20_contract_address = ) AND (`from` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155 OR `to` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155) ORDER BY block_number DESC,id DESC LIMIT 3 OFFSET 0",
				Response: []map[string]interface{}{
					{
						"block_number":           1665113,
						"hash":                   "0x691c5ff7af12681d514dcce7d37f0be1328f647597f175d3b55ecd09fd0965c9",
						"nonce":                  "16829",
						"block_hash":             "0xc47ae8c594ce11bff788a84641dfacec04b82d73b8adc269a1da15e5cb991866",
						"from":                   "0xcda3283e436807a798fc2f8b07be20f2b3b4d155",
						"to":                     "",
						"value":                  "0",
						"gas":                    "500000",
						"gas_price":              "100000000000",
						"timestamp":              1505289966,
						"transaction_index":      "0",
						"is_contract":            false,
						"status":                 0,
						"confirmed_num":          0,
						"erc20_contract_address": "",
						"url": "https://etherscan.io/tx/0x691c5ff7af12681d514dcce7d37f0be1328f647597f175d3b55ecd09fd0965c9",
					},
					{
						"block_number":           1665077,
						"hash":                   "0xa5c45b210002354d1243de8e8ccd07638c5e0b866c022b5f67d8daf6bfb9b33b",
						"nonce":                  "16828",
						"block_hash":             "0x5d158f977399c94ea7c068211e32355d604c0df6d0af66b58048d525c992426b",
						"from":                   "0xcda3283e436807a798fc2f8b07be20f2b3b4d155",
						"to":                     "0xad5c52d34b5e4cf0c030ff5726ba44b2c4802e22",
						"value":                  "0",
						"gas":                    "500000",
						"gas_price":              "100000000000",
						"timestamp":              1505289370,
						"transaction_index":      "0",
						"is_contract":            true,
						"status":                 0,
						"confirmed_num":          0,
						"erc20_contract_address": "",
						"url": "https://etherscan.io/tx/0xa5c45b210002354d1243de8e8ccd07638c5e0b866c022b5f67d8daf6bfb9b33b",
					},
					{
						"block_number":           1665074,
						"hash":                   "0x13e7b40ed70433e607dd7a7a41f9e77e4f5e5231acafb345d23fcf5bb4e3562c",
						"nonce":                  "16827",
						"block_hash":             "0x297345dd4c83b622d9ad0b8e53d6dabd3814b5445d3ef2c5f7522e473e862338",
						"from":                   "0xcda3283e436807a798fc2f8b07be20f2b3b4d155",
						"to":                     "",
						"value":                  "0",
						"gas":                    "500000",
						"gas_price":              "100000000000",
						"timestamp":              1505289340,
						"transaction_index":      "0",
						"is_contract":            false,
						"status":                 0,
						"confirmed_num":          0,
						"erc20_contract_address": "",
						"url": "https://etherscan.io/tx/0x13e7b40ed70433e607dd7a7a41f9e77e4f5e5231acafb345d23fcf5bb4e3562c",
					},
				},
			},
		})

		txList, total := EthereumTransactionList("0xcda3283e436807a798fc2f8b07be20f2b3b4d155", "", 1, 3)
		if total != 6 {
			t.Errorf("count of tx does not match")
		}
		for _, tx := range txList {
			if tx.ERC20ContractAddress != "" ||
				(tx.From != "0xcda3283e436807a798fc2f8b07be20f2b3b4d155" &&
					tx.To != "0xcda3283e436807a798fc2f8b07be20f2b3b4d155") {
				t.Errorf("tx does not match")
			}
		}

		txList, total = EthereumTransactionList("", "", 1, 4)
		if len(txList) != 0 || total != 0 {
			t.Errorf("txList is not empty or count of tx is not zero")
		}

		txList, total = EthereumTransactionList("0xcda3283e436807a798fc2f8b07be20f2b3b4d155", "", 0, 0)
		if len(txList) != 0 || total != 0 {
			t.Errorf("txList is not empty or count of tx is not zero")
		}

	})

	t.Run("erc20 tx list", func(t *testing.T) {
		mocket.Catcher.Reset()
		mocket.Catcher.Attach([]*mocket.FakeResponse{
			{
				Pattern:  "SELECT count(*) FROM \"transaction_infos\"  WHERE (erc20_contract_address = 0x6f72a579ec73438ceb70fcc5c00e7e8579e5b9ee) AND (`from` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155 OR `to` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155)",
				Response: []map[string]interface{}{{"total": 3}},
			},
			{
				Pattern: "SELECT * FROM \"transaction_infos\"  WHERE (erc20_contract_address = 0x6f72a579ec73438ceb70fcc5c00e7e8579e5b9ee) AND (`from` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155 OR `to` = 0xcda3283e436807a798fc2f8b07be20f2b3b4d155) ORDER BY block_number DESC,id DESC LIMIT 4 OFFSET 0",
				Response: []map[string]interface{}{
					{
						"block_number":           1429526,
						"hash":                   "0xc62467409f97106402cf8f07d44e067fc9e339eda24f5ec604da109193a9d53b",
						"nonce":                  "67",
						"block_hash":             "0xb5332f6fa082a78e827e6978ce76d6fe9a00488a7859bec8eb8f9ebdea48ce5a",
						"from":                   "0xcda3283e436807a798fc2f8b07be20f2b3b4d155",
						"to":                     "0xf81d36ac26f173c32ad2c43bc26083b46ac769e9",
						"value":                  "1000000000000000000",
						"gas":                    "90000",
						"gas_price":              "30000000000",
						"timestamp":              1501898677,
						"transaction_index":      "0",
						"is_contract":            true,
						"status":                 0,
						"confirmed_num":          12,
						"erc20_contract_address": "0x6f72a579ec73438ceb70fcc5c00e7e8579e5b9ee",
						"url": "https://etherscan.io/tx/0xc62467409f97106402cf8f07d44e067fc9e339eda24f5ec604da109193a9d53b",
					},
					{
						"block_number":           1429515,
						"hash":                   "0xea6282ade914708b79fd40746e6d22d795f6484b0a6416e9b8daef8739bd44ce",
						"nonce":                  "66",
						"block_hash":             "0xf15b62e0659a8c09de7220b2f3a39935f919a8aa24bd345e47e99df65273d646",
						"from":                   "0xcda3283e436807a798fc2f8b07be20f2b3b4d155",
						"to":                     "0xae270efa70b943b4f4b84009a9f6b80f40414194",
						"value":                  "1000000000000000000",
						"gas":                    "90000",
						"gas_price":              "20000000000",
						"timestamp":              1501898509,
						"transaction_index":      "4",
						"is_contract":            true,
						"status":                 0,
						"confirmed_num":          12,
						"erc20_contract_address": "0x6f72a579ec73438ceb70fcc5c00e7e8579e5b9ee",
						"url": "https://etherscan.io/tx/0xea6282ade914708b79fd40746e6d22d795f6484b0a6416e9b8daef8739bd44ce",
					},
				},
			},
		})

		txList, total := EthereumTransactionList("0xCDA3283e436807A798fc2f8b07be20f2b3b4d155", "0x6f72a579ec73438ceb70FCC5c00e7e8579e5b9EE", 1, 4)
		if total != 3 {
			t.Errorf("count of tx does not match")
		}
		for _, tx := range txList {
			if tx.ERC20ContractAddress != "0x6f72a579ec73438ceb70fcc5c00e7e8579e5b9ee" ||
				(tx.From != "0xcda3283e436807a798fc2f8b07be20f2b3b4d155" &&
					tx.To != "0xcda3283E436807a798fc2f8b07be20f2b3b4d155") {
				t.Errorf("tx does not match")
			}
		}
	})

}
