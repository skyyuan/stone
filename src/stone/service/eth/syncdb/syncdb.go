package syncdb

import (
	"encoding/json"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"stone/common"
	"stone/common/ethcomm"
	"stone/nsqs"
	"stone/service/eth"
	"stone/service/utils"
)

// LastBlockNumber synced last block number
var LastBlockNumber int64

// AutoFindErc20 sync erc20 flag
var DisableAutoErc20 bool

var jobDone = true

// UnitRepairBlockNum UnitRepairBlockNum
const UnitRepairBlockNum = int64(50000)
const connectionNum = 10

type repairBlockInfo struct {
	from int64
	to   int64
}

var needRepairBlock []repairBlockInfo

// MigrateDB migrate db
func MigrateDB() {
	db := ethcomm.DBBegin()
	db.AutoMigrate(&eth.TransactionInfo{}, &eth.BlockInfo{})
	db.DBCommit()
}

// CheckAndRepairBlockData CheckAndRepairBlockData
func CheckAndRepairBlockData() {
	currentBlockNum, err := getCurrentBlockNumber()
	if err != nil {
		common.Logger.Error(err)
		return
	}
	// set Global LastBlockNumber
	LastBlockNumber = currentBlockNum
	var count int64
	updateLast12Blocks(currentBlockNum)
	go updateUnconfirmedBlocks(currentBlockNum)

	db := ethcomm.DBBegin()
	defer db.DBRollback()

	db.Model(&eth.BlockInfo{}).Where("is_confirmed = ?", true).Count(&count)

	if count == currentBlockNum-eth.ConfirmedNum {
		return
	}
	needRepairBlock = make([]repairBlockInfo, 0, 0)
	syncedBlockNum := currentBlockNum - eth.ConfirmedNum

	findNeedRepairBlockInfos(db, 0, syncedBlockNum)
	db.DBRollback()

	unitBlockGorutine := len(needRepairBlock) / connectionNum
	for idx := 0; idx < 10; idx++ {
		go repairBlockInAGorutine(needRepairBlock[idx*unitBlockGorutine:(idx+1)*unitBlockGorutine], currentBlockNum)
	}
	modNum := len(needRepairBlock) % connectionNum
	go repairBlockInAGorutine(needRepairBlock[len(needRepairBlock)-modNum:], currentBlockNum)
}

func updateLast12Blocks(currentBlockNum int64) {
	for idx := int64(0); idx < eth.ConfirmedNum; idx++ {
		updateUnconfirmedBlock(currentBlockNum-idx, currentBlockNum)
	}
}

func updateUnconfirmedBlocks(currentBlockNum int64) {
	var blockInfos []*eth.BlockInfo
	db := ethcomm.DBBegin()

	db.Model(&eth.BlockInfo{}).Where("is_confirmed = ?", false).Find(&blockInfos)
	db.DBRollback()

	var wg sync.WaitGroup
	for _, item := range blockInfos {
		wg.Add(1)
		go func(blockNum, currentBlockNum int64) {
			defer wg.Done()
			updateUnconfirmedBlock(blockNum, currentBlockNum)
		}(item.Number, currentBlockNum)
	}
	wg.Wait()
}

func repairBlockInAGorutine(repairBlockInfos []repairBlockInfo, currentBlockNum int64) {
	for _, item := range repairBlockInfos {
		syncConfirmedBlocks(item.from, item.to, currentBlockNum)
	}
}

func findNeedRepairBlockInfos(db *ethcomm.GormDB, from int64, to int64) {
	if to-from < UnitRepairBlockNum {
		needRepairBlock = append(needRepairBlock, repairBlockInfo{from: from, to: to})
		return
	}
	med := (from + to) / 2
	leftCnt := int64(0)
	db.Model(&eth.BlockInfo{}).Where("number >= ? and number <= ?", from, med).Count(&leftCnt)
	if leftCnt != med-from+1 {
		findNeedRepairBlockInfos(db, from, med)
	}
	rightCnt := int64(0)
	db.Model(&eth.BlockInfo{}).Where("number >= ? and number <= ?", med+1, to).Count(&rightCnt)
	if rightCnt != to-med {
		findNeedRepairBlockInfos(db, med+1, to)
	}
}

func syncConfirmedBlocks(fromBlockNum int64, toBlockNum int64, currentBlockNum int64) {
	for idx := fromBlockNum; idx <= toBlockNum; idx++ {
		syncConfirmedBlock(idx, currentBlockNum)
	}
}

// SyncEthDB SyncEthDB
func SyncEthDB() {

	if !jobDone {
		common.Logger.Error("jobDone is false at: ", LastBlockNumber)
		return
	}
	jobDone = false
	defer func() {
		jobDone = true
	}()
	currentBlockNum, err := getCurrentBlockNumber()
	common.Logger.Info("sync block: ", currentBlockNum)
	if err != nil {
		common.Logger.Error(err)
		return
	}
	var wg sync.WaitGroup
	times := currentBlockNum - LastBlockNumber
	for idx := int64(0); idx < times; idx++ {
		tmpBlockNum := LastBlockNumber + idx + 1
		common.Logger.Info("syncing :", tmpBlockNum)
		wg.Add(1)
		go func(tmpBlockNum, currentBlockNum int64) {
			defer wg.Done()
			syncNewBlock(tmpBlockNum, currentBlockNum)
		}(tmpBlockNum, currentBlockNum)

	}
	LastBlockNumber = currentBlockNum

	updateUnconfirmedBlocks(currentBlockNum)
	wg.Wait()
	common.Logger.Info("lastBlockNumber is sync: ", LastBlockNumber)
}

func getCurrentBlockNumber() (num int64, err error) {
	currentBlockNumStr, err := eth.BlockNumber()
	if err != nil {
		return
	}
	num, err = utils.HexToInt(currentBlockNumStr)
	return
}

func syncConfirmedBlock(blockNum int64, currentBlockNum int64) error {
	common.Logger.Info("sync confirmed block:", blockNum)
	db := ethcomm.DBBegin()
	defer db.DBRollback()
	tmpBlock := eth.BlockInfo{}
	notFound := db.Where("number = ?", blockNum).Find(&tmpBlock).RecordNotFound()
	if !notFound {
		return nil
	}
	err := createBlockToDB(db, blockNum, currentBlockNum)
	db.DBCommit()
	return err
}

func syncNewBlock(blockNum int64, currentBlockNum int64) (err error) {
	db := ethcomm.DBBegin()
	defer db.DBRollback()
	tmpBlock := &eth.BlockInfo{}
	notFound := db.Where("number = ?", blockNum).Find(tmpBlock).RecordNotFound()

	if notFound {
		err = createBlockToDB(db, blockNum, currentBlockNum)
	} else {
		err = updateBlockToDB(db, blockNum, currentBlockNum)
	}
	db.DBCommit()

	if !DisableAutoErc20 {
		err = nsqs.PostTopic("eth_sync_block_erc20", blockNum)
	}
	if err != nil {
		common.Logger.Error(err)
	}
	return err
}

func createBlockToDB(db *ethcomm.GormDB, blockNum int64, currentBlockNum int64) error {
	blockInfo, err := eth.GetBlockInfo(utils.Int64ToHex(blockNum))
	if err != nil {
		return err
	}
	txns := make([]map[string]interface{}, 0)
	err = json.Unmarshal(blockInfo.Transactions, &txns)
	if err != nil {
		return err
	}
	confirmedNum := currentBlockNum - blockInfo.Number
	if confirmedNum > eth.ConfirmedNum {
		confirmedNum = eth.ConfirmedNum
	}

	txnLen := len(txns)
	concurrencyNum := 20
	jobs := make(chan *eth.TransactionInfo, txnLen)
	results := make(chan *txnWorkerResult)

	for _, item := range txns {
		tmpTxn := &eth.TransactionInfo{}
		eth.MapAdapter(item, tmpTxn)
		jobs <- tmpTxn
	}

	close(jobs)

	for i := 0; i < concurrencyNum; i++ {
		go txnCreateWorker(confirmedNum, blockInfo.Timestamp, false, jobs, results)
	}

	for idx := 0; idx < txnLen; idx++ {
		workerRes := <-results
		if workerRes.err != nil {
			continue
		}
		if workerRes.txn.To == "0x0000000000000000000000000000000000000020" {
			continue
		}
		common.Logger.Info("txnCnt=", idx, " txHash=", workerRes.txn.Hash)
		err = db.Create(workerRes.txn).Error
		if err != nil {
			common.Logger.Error(err)
		}
		createERC20TxnFromLogs(db, workerRes.txn)
	}

	if confirmedNum >= eth.ConfirmedNum {
		blockInfo.IsConfirmed = true
	} else {
		blockInfo.IsConfirmed = false
	}
	err = db.Create(blockInfo).Error
	if err != nil {
		common.Logger.Error(err)
	}
	return nil
}

func updateBlockToDB(db *ethcomm.GormDB, blockNum int64, currentBlockNum int64) error {
	err := db.Where("block_number = ?", blockNum).Delete(eth.TransactionInfo{}).Error
	if err != nil {
		common.Logger.Error(err)
	}
	err = db.Where("number = ?", blockNum).Delete(eth.BlockInfo{}).Error
	if err != nil {
		common.Logger.Error(err)
	}
	return createBlockToDB(db, blockNum, currentBlockNum)
}

func updateUnconfirmedBlock(blockNum int64, currentBlockNum int64) error {

	db := ethcomm.DBBegin()
	defer db.DBRollback()

	blockInfo, err := eth.GetBlockInfo(utils.Int64ToHex(blockNum))
	if err != nil {
		return err
	}
	tmpBlockInfo := &eth.BlockInfo{}
	notFound := db.Where("hash = ?", blockInfo.Hash).Find(tmpBlockInfo).RecordNotFound()
	if notFound {
		updateBlockToDB(db, blockNum, currentBlockNum)
	} else {
		confirmedNum := currentBlockNum - blockNum
		if confirmedNum > eth.ConfirmedNum {
			confirmedNum = eth.ConfirmedNum
		}
		//batch update confirmedNum
		err = db.Model(&eth.TransactionInfo{}).Where("block_number = ?", blockInfo.Number).Update("confirmed_num", confirmedNum).Error
		if err != nil {
			common.Logger.Error(err)
		}
		if confirmedNum >= eth.ConfirmedNum {
			err = db.Model(tmpBlockInfo).Update("is_confirmed", true).Error
			if err != nil {
				common.Logger.Error(err)
			}
		}
	}
	db.DBCommit()
	return nil
}

func isErc20TxnEvent(topics []ethcommon.Hash) bool {
	return len(topics) == 3 && eth.EventTransferHash == topics[0].String()
}

func createERC20TxnFromLogs(db *ethcomm.GormDB, txn *eth.TransactionInfo) {
	receiptLogs := []types.Log{}
	err := json.Unmarshal(txn.Receipt.Logs, &receiptLogs)
	if err != nil {
		common.Logger.Error(err)
		return
	}
	for _, item := range receiptLogs {
		if isErc20TxnEvent(item.Topics) {
			createERC20TxnFromLog(db, txn, item)
		}
	}
}

func createERC20TxnFromLog(db *ethcomm.GormDB, txn *eth.TransactionInfo, log types.Log) {
	newTxn := *txn
	newTxn.ID = 0
	newTxn.From = eth.HexPrefix + log.Topics[1].String()[eth.AddressInHashIndex:]
	newTxn.To = eth.HexPrefix + log.Topics[2].String()[eth.AddressInHashIndex:]
	newTxn.ERC20ContractAddress = log.Address.Hex()
	value, _ := utils.ParseBig256((ethcommon.BytesToHash(log.Data)).String())
	newTxn.Value = value.Text(10)
	err := db.Create(&newTxn).Error
	if err != nil {
		common.Logger.Error(err)
	}
}

type txnWorkerResult struct {
	txn *eth.TransactionInfo
	err error
}

func txnCreateWorker(confirmedNum, timestamp int64, containERC20Flag bool, jobs <-chan *eth.TransactionInfo, results chan<- *txnWorkerResult) {
	for txn := range jobs {
		receipt, err := eth.TransactionReceipt(txn.Hash)
		if err != nil {
			common.Logger.Error(txn.Hash, err)
			results <- &txnWorkerResult{nil, err}
			continue
		}
		txn.Receipt = receipt
		txn.Status, err = eth.TxnStatus(txn, receipt)
		if err != nil {
			common.Logger.Error(txn.Hash, err)
			results <- &txnWorkerResult{nil, err}
			continue
		}
		if txn.To != "" {
			txn.IsContract, err = eth.IsContract(txn.To)
		}
		if err != nil {
			common.Logger.Error(txn.Hash, err)
			results <- &txnWorkerResult{nil, err}
			continue
		}
		txn.ConfirmedNum = confirmedNum
		txn.Timestamp = timestamp
		if containERC20Flag && eth.IsERC20Txn(txn) {
			eth.GenERC20Txn(txn)
		}
		results <- &txnWorkerResult{txn, nil}
	}
}
