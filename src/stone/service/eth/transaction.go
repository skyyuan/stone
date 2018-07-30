package eth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	ecom "github.com/ethereum/go-ethereum/common"
	"stone/common"
	"stone/service/etherscan"
	"stone/service/utils"

	"stone/common/ethcomm"
)

// DetailTransactionInfo transaction detail info
func DetailTransactionInfo(txn string) (*TransactionInfo, error) {
	var resp TransactionInfo
	var err error
	rcv := make(map[string]interface{})
	err = endpointsManager.RPC(&rcv, "eth_getTransactionByHash", txn)
	if err != nil {
		common.Logger.Debug(err)
		return nil, err
	}

	// get block timestamp
	var blockinfo map[string]interface{}
	err = endpointsManager.RPC(&blockinfo, "eth_getBlockByHash", rcv["blockHash"], false)
	if err != nil {
		common.Logger.Debug(err)
		return nil, err
	}
	rcv["timestamp"] = blockinfo["timestamp"]
	MapAdapter(rcv, &resp)

	// update status
	receipt, err := TransactionReceipt(resp.Hash)
	if err != nil {
		common.Logger.Debug(err)
		return nil, err
	}
	status, err := TxnStatus(&resp, receipt)
	if err != nil {
		common.Logger.Debug(err)
		return nil, err
	}
	resp.Status = status

	v, _ := utils.HexToInt(receipt.GasUsed)
	resp.GasUsed = strconv.Itoa(int(v))

	// check toAddress is contract address
	if resp.To != "" {
		isContract, err := IsContract(resp.To)
		if err != nil {
			common.Logger.Debug(err)
			return nil, err
		}
		resp.IsContract = isContract
	}

	// check txn is confirmed
	confirmedNum, err := GetComfirmedNum(&resp)
	resp.ConfirmedNum = confirmedNum
	resp.URL = GetURL(resp.Hash)
	if len(resp.Input) > 2 { //IsERC20Txn(&resp) {
		//GenERC20Txn(&resp)
		if v, ok := ATMFunction[resp.Input[2:10]]; ok {
			para := make([]string, (len(resp.Input)-10)/64+1)
			para[0] = v
			fmt.Println(para[0])
			k := 1
			i := 10
			for i < len(resp.Input) {
				para[k] = resp.Input[i : i+64]
				fmt.Println(para[k])
				k = k + 1
				i = i + 64
			}
			resp.InputData = para
		} else {
			fmt.Println("Key Not Found:", resp.Input[2:10])
		}
		db := ethcomm.DBBegin()
		defer db.DBRollback()
		total := 0
		txlist := make([]*TransactionInfo, 0)
		allTxs := db.Model(&TransactionInfo{}).
			Where("erc20_contract_address != ''").
			Where("`hash` = ?", resp.Hash)
		allTxs.Count(&total)
		allTxs.Order("block_number DESC").Order("id DESC").
			Find(&txlist)
		internal := make([]string, total)
		for index, item := range txlist {
			internal[index] = item.From + "," + item.To + "," + item.Value
		}
		resp.InternalTransaction = internal
	}

	return &resp, err
}

// TransactionReceipt transaction receipt
func TransactionReceipt(txnHash string) (*TransactionReceiptInfo, error) {
	var resp TransactionReceiptInfo
	err := endpointsManager.RPC(&resp, "eth_getTransactionReceipt", txnHash)
	if err != nil {
		common.Logger.Debug(err)
		return nil, err
	}
	return &resp, err
}

// GetBlockInfo get block info with the block number
func GetBlockInfo(blockNum string) (*BlockInfo, error) {
	// get block timestamp
	var blockinfo BlockInfo
	rcv := make(map[string]interface{})

	err := endpointsManager.RPC(&rcv, "eth_getBlockByNumber", blockNum, true)
	if err != nil {
		common.Logger.Debug(err)
		return nil, err
	}
	MapAdapter(rcv, &blockinfo)
	if rcv["transactions"] == nil {
		blockinfo.Transactions = []byte{}
	} else {
		txnStr, err := json.Marshal(rcv["transactions"].([]interface{}))
		if err != nil {
			return nil, err
		}
		blockinfo.Transactions = []byte(txnStr)
	}
	var sealer string
	err = endpointsManager.RPC(&sealer, "clique_getSealer", blockNum)
	if err != nil {
		common.Logger.Error(blockNum," get sealer err: ", err)
	}else{
		blockinfo.Miner = sealer
	}
	return &blockinfo, nil
}

// SendRawTransaction sends transaction
func SendRawTransaction(hexstr string) (string, error) {
	var resp string
	go etherscan.SendRawTransaction(isTestnet, hexstr)
	err := endpointsManager.RPC(&resp, "eth_sendRawTransaction", hexstr)
	if err != nil {
		common.Logger.Debug(err)
		return "", err
	}
	return resp, nil

}

//EthereumTransactionList ethereum transaction list
func EthereumTransactionList(userAddress, contractAddress string,
	currentPage, perPage int) ([]*TransactionInfo, int) {

	total := 0
	txlist := make([]*TransactionInfo, 0)
	if currentPage > 0 && perPage > 0 {
		fromIdx := (currentPage - 1) * perPage

		db := ethcomm.DBBegin()
		defer db.DBRollback()

		if ecom.IsHexAddress(userAddress) {
			userAddress = strings.ToLower(userAddress)
			if ecom.IsHexAddress(contractAddress) {
				contractAddress = strings.ToLower(contractAddress)
				allTxs := db.Model(&TransactionInfo{}).
					Where("erc20_contract_address = ?", contractAddress).
					Where("`from` = ? OR `to` = ?", userAddress, userAddress)
				allTxs.Count(&total)
				allTxs.Order("block_number DESC").Order("id DESC").Offset(fromIdx).
					Limit(perPage).Find(&txlist)
			} else {
				allTxs := db.Model(&TransactionInfo{}).
					Where("erc20_contract_address = ''").
					Where("`from` = ? OR `to` = ?", userAddress, userAddress)
				allTxs.Count(&total)
				allTxs.Order("block_number DESC").Order("id DESC").Offset(fromIdx).
					Limit(perPage).Find(&txlist)
			}

		} else {
			allTxs := db.Model(&TransactionInfo{}).
				Where("erc20_contract_address = ''")
			allTxs.Count(&total)
			allTxs.Order("block_number DESC").Order("id DESC").Offset(fromIdx).
				Limit(perPage).Find(&txlist)
		}

		for _, item := range txlist {
			item.URL = GetURL(item.Hash)
		}
	}
	return txlist, total
}

//EthereumBlockList ethereum transaction list
func EthereumBlockList(currentPage, perPage int) ([]*BlockInfo, int) {

	total := 0
	blocklist := make([]*BlockInfo, 0)
	if currentPage > 0 && perPage > 0 {
		fromIdx := (currentPage - 1) * perPage

		db := ethcomm.DBBegin()
		defer db.DBRollback()
		allBlocks := db.Model(&BlockInfo{})
		allBlocks.Count(&total)
		allBlocks.Order("number DESC").Order("id DESC").Offset(fromIdx).
			Limit(perPage).Find(&blocklist)

		for _, item := range blocklist {
			item.TxnCounts = EthereumTransactionCount(strconv.Itoa(int(item.Number)))
		}

	}
	return blocklist, total
}

//EthereumTransactionCount ethereum transaction count in block_number
func EthereumTransactionCount(block_number string) int {
	total := 0
	db := ethcomm.DBBegin()
	defer db.DBRollback()

	allTxs := db.Model(&TransactionInfo{}).
		Where("`block_number` = ?", block_number).
		Where("erc20_contract_address = ''")
	allTxs.Count(&total)
	return total
}

//EthereumTransactionListByBlock ethereum transaction list in block_number
func EthereumTransactionListByBlock(block_number string,
	currentPage, perPage int) ([]*TransactionInfo, int) {
	total := 0
	txlist := make([]*TransactionInfo, 0)
	if currentPage > 0 && perPage > 0 {
		fromIdx := (currentPage - 1) * perPage

		db := ethcomm.DBBegin()
		defer db.DBRollback()

		allTxs := db.Model(&TransactionInfo{}).
			Where("erc20_contract_address = ''").
			Where("`block_number` = ?", block_number)
		allTxs.Count(&total)
		allTxs.Order("id DESC").Offset(fromIdx).
			Limit(perPage).Find(&txlist)

		for _, item := range txlist {
			item.URL = GetURL(item.Hash)
		}
	}
	return txlist, total
}

//EthereumBlock ethereum block
func EthereumBlock(blockNumber string) *BlockInfo {

	block := &BlockInfo{}

	db := ethcomm.DBBegin()
	defer db.DBRollback()
	allBlocks := db.Model(&BlockInfo{}).
		Where("`number` = ?", blockNumber)
	allBlocks.Find(&block)
	block.TxnCounts = EthereumTransactionCount(blockNumber)
	return block
}

// DetailTransactionInfoFromDB get tx from db with txid
func DetailTransactionInfoFromDB(txid uint) (*TransactionInfo, error) {
	db := ethcomm.DBBegin()
	defer db.DBRollback()
	txInfo := &TransactionInfo{}
	txInfo.ID = txid
	err := db.Find(txInfo).Error
	txInfo.URL = GetURL(txInfo.Hash)
	return txInfo, err
}

// HttpGetATMPrice get ATM price
func HttpGetATMPrice() (*ATMPrice, error) {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/attention-token-of-media/?convert=CNY")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result := &ATMPrice{}
	result.PriceCny = string(body)
	//fmt.Println(result.PriceCny)
	return result, err
}
