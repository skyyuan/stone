package eth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"stone/common"
	"stone/common/ethcomm"
	"stone/service/utils"
)

var (
	gethEndpoint     string
	isTestnet        bool
	endpointsManager *EndpointsManager
)

// BlockInfo block info
type BlockInfo struct {
	common.Model
	Number       int64  `gorm:"unique_index" conv:"int64"`
	Hash         string `gorm:"type:varchar(255);unique_index" conv:""`
	ParentHash   string `gorm:"type:varchar(255)" conv:""`
	Timestamp    int64  `conv:"int64"`
	Transactions []byte `gorm:"-"`
	IsConfirmed  bool   `gorm:"index"`
	Difficulty   int64  `conv:"int64"`
	ExtraData    string `gorm:"type:varchar(255)" conv:""`
	GasLimit     int64  `conv:"int64"`
	GasUsed      int64  `conv:"int64"`
	//LogsBloom        string `gorm:"type:varchar(255)" conv:""`
	Miner            string `gorm:"type:varchar(255)" conv:""`
	MixHash          string `gorm:"type:varchar(255)" conv:""`
	Nonce            string `gorm:"type:varchar(255)" conv:""`
	ReceiptsRoot     string `gorm:"type:varchar(255)" conv:""`
	Sha3Uncles       string `gorm:"type:varchar(255)" conv:""`
	Size             int64  `conv:"int64"`
	StateRoot        string `gorm:"type:varchar(255)" conv:""`
	TotalDifficulty  int64  `conv:"int64"`
	TransactionsRoot string `gorm:"type:varchar(255)" conv:""`
	TxnCounts        int    `json:"TxnCounts" gorm:"-"`
}

// EthereumNodeInfo represents an ethereum transaction hash
type EthereumNodeInfo struct {
	URL       string `json:"url"`
	Version   string `json:"version"`
	Is_mining bool   `json:"is_mining"`
	Miner     string `json:"miner"`
	Is_alive  bool   `json:"is_alive"`
}

// TransactionInfo transaction info
type TransactionInfo struct {
	common.Model
	BlockNumber          int64                   `json:"blockNumber" conv:"int64" gorm:"index:idx_txninfo_blocknumber"`
	Hash                 string                  `json:"hash" conv:"" gorm:"index:idx_txninfo_hash"`
	Nonce                string                  `json:"nonce" conv:"int256str"`
	BlockHash            string                  `json:"blockHash" conv:"" gorm:"index:idx_txninfo_blockhash"`
	From                 string                  `json:"from" conv:"" gorm:"index:idx_txninfo_from"`
	To                   string                  `json:"to" conv:"" gorm:"index:idx_txninfo_to"`
	Value                string                  `json:"value" conv:"int256str"`
	Input                string                  `gorm:"type:text" json:"input" conv:""`
	InputData            []string                `gorm:"-" json:"inputData" conv:""`
	InternalTransaction  []string                `gorm:"-" json:"internalTransaction" conv:""`
	Gas                  string                  `json:"gas" conv:"int256str"`
	GasPrice             string                  `json:"gasPrice" conv:"int256str"`
	Timestamp            int64                   `json:"timestamp" conv:"int64"`
	TransactionIndex     string                  `json:"transactionIndex" conv:"int256str"`
	IsContract           bool                    `json:"isContract"`
	Status               int                     `json:"status"`
	ConfirmedNum         int64                   `json:"confirmedNum" gorm:"index:idx_txninfo_confirmed_num"`
	ERC20ContractAddress string                  `json:"erc20ContractAddress" gorm:"column:erc20_contract_address;index:idx_txninfo_erc20_contract_address"`
	URL                  string                  `json:"url" gorm:"-"`
	Receipt              *TransactionReceiptInfo `json:"-" gorm:"-"`
	GasUsed              string                  `json:"gasUsed" gorm:"-"`
}

// TransactionReceiptInfo transaction receipt info
type TransactionReceiptInfo struct {
	TransactionHash string          `json:"transactionHash"`
	BlockHash       string          `json:"blockHash"`
	BlockNumber     string          `json:"blockNumber"`
	GasUsed         string          `json:"gasUsed"`
	Logs            json.RawMessage `json:"logs"`
}
type ATMPrice struct {
	common.Model
	PriceCny string `json:"priceCny"`
}

// BalanceAmount balance amount
type BalanceAmount struct {
	Amount string `json:"amount"`
}

// GasPrice represents gasPrice
type GasPrice struct {
	GasPrice string `json:"gasPrice"`
}

// ServiceInit init ethereum service. getEndpoint is geth url
func ServiceInit(config ethcomm.Config) {
	isTestnet = config.EthTestNet()

	endpointsManager = NewEndPointsManager()

	file, err := os.Open(config.GethEndpoint())
	defer file.Close()
	if err != nil {
		fmt.Println("ServiceInit err=", err)
		return
	}
	bufferedReader := bufio.NewReader(file)

	for {
		if dataString, err := bufferedReader.ReadString('\n'); err == nil && dataString[0:4] == "http" {
			fmt.Println("endpoint=", dataString[0:len(dataString)-1])
			endpointsManager.AddEndPoint(dataString[0:len(dataString)-1], 1)
		} else {
			break
		}
	}
	go endpointsManager.Run()
}

// ServiceDone close endpointsManager
func ServiceDone() {
	endpointsManager.Stop()
}

// MapAdapter adapt m to i.
func MapAdapter(m map[string]interface{}, i interface{}) {
	iMeta := reflect.ValueOf(i).Elem()
	for i := 0; i < iMeta.NumField(); i++ {
		valueField := iMeta.Field(i)
		if !valueField.CanSet() {
			continue
		}
		typeField := iMeta.Type().Field(i)
		typeName := typeField.Name
		v, ok := m[strings.ToLower(typeName[:1])+typeName[1:]]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}

		tags := typeField.Tag
		conv := tags.Get("conv")
		switch conv {
		case "int64":
			i64, err := utils.HexToInt(s)
			if err == nil {
				valueField.SetInt(i64)
			}
		case "int256str":
			i256, ok := utils.ParseBig256(s)
			if ok {
				valueField.SetString(i256.Text(10))
			}
		case "":
			valueField.SetString(s)
		default:
			continue
		}
	}
}

// GetURL get url
func GetURL(hash string) string {
	if isTestnet {
		return RopstenURL + hash
	}
	return EtherscanURL + hash
}
