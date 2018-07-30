package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGasPrice(t *testing.T) {
	gasPrice, err := GetGasPrice()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, gasPrice, nil)
}

func TestGetTransactionCount(t *testing.T) {
	count, err := GetTransactionCount("0x413b2b64fd25ab1b77d117005285dd4d40c509f7")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, count, "1", "The two value should be the same.")

	count, err = GetTransactionCount("0x413b2b64fd25ab1b77d117005285dd4d40c509f8")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, count, "0", "tx count for non-exsting account should be 0.")
}

func TestGetCode(t *testing.T) {
	code, err := GetCode("0x54fb2b073926f20aa46604b00ccab89e50b5505f")
	if err != nil {
		t.Fatal(err)
	}
	assert.Regexp(t, "^0[xX][0-9a-fA-F]+$", code, "contract code should be found")

	code, err = GetCode("0x33853a41d81187759ba54fbc1492120bebec3d6b")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotRegexp(t, "^0[xX][0-9a-fA-F]+$", code, "contract code should not be found")
}

func TestStatus(t *testing.T) {
	hash := "0xd2397b1e1a7b3230616e6232c2de205e1cf4d2870c7204f979a79f8ed11a9a90"
	txninfo, err := DetailTransactionInfo(hash)
	if err != nil {
		t.Fatal("Error", err)
	}
	receipt, _ := TransactionReceipt(txninfo.Hash)
	status, err := TxnStatus(txninfo, receipt)
	if err != nil {
		t.Fatal("Error", err)
	}
	assert.Equal(t, status, 0, "Status does not match")

	// Non-existing tx hash makes servide panic
}

func TestGetComfirmedNum(t *testing.T) {
	hash := "0x3965d599457c146e28979eb783c4cf9ed8d1e42c7d56a89d6a55b3dfb44157f0"
	txninfo, err := DetailTransactionInfo(hash)
	if err != nil {
		t.Fatal(err)
	}
	confirmedNum, err := GetComfirmedNum(txninfo)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, confirmedNum, int64(12), "sss")
	txninfo = &TransactionInfo{Status: 2}
	confirmedNum, err = GetComfirmedNum(txninfo)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, confirmedNum, int64(0), "sss")

}

func TestIsERC20Txn(t *testing.T) {
	txninfo := &TransactionInfo{Input: "0xa9059cbb000000000000000000000000bdd7ebf5709fa00ea277ae1458ea6d64f26652ac000000000000000000000000000000000000000000000001a055690d9db80000"}
	isErc20 := IsERC20Txn(txninfo)
	assert.Equal(t, isErc20, true, "should be startwith 0xa9059cbb")

	txninfo = &TransactionInfo{Input: "0x"}
	assert.False(t, IsERC20Txn(txninfo))
	txninfo = &TransactionInfo{Input: ""}
	assert.False(t, IsERC20Txn(txninfo))
	txninfo = &TransactionInfo{Input: "0xsdjfls"}
	assert.False(t, IsERC20Txn(txninfo))

}

func TestGenERC20Txn(t *testing.T) {
	txninfo := &TransactionInfo{To: "0x81BfB6A2Db736c5EC06DdF4654478CF78B3E0bE7", Input: "0xa9059cbb000000000000000000000000bdd7ebf5709fa00ea277ae1458ea6d64f26652ac000000000000000000000000000000000000000000000001a055690d9db80000"}
	GenERC20Txn(txninfo)

	assert.Equal(t, txninfo.To, "0xbdd7ebf5709fa00ea277ae1458ea6d64f26652ac", "to address bdd7ebf5709fa00ea277ae1458ea6d64f26652ac")
	assert.Equal(t, txninfo.Value, "30000000000000000000", "value")
	assert.Equal(t, txninfo.ERC20ContractAddress, "0x81BfB6A2Db736c5EC06DdF4654478CF78B3E0bE7", "ERC20ContractAddress")
}

func TestIsContract(t *testing.T) {
	isContract, err := IsContract("0x54fb2b073926f20aa46604b00ccab89e50b5505f")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, isContract, true, "IsContract fails")

	isContract, err = IsContract("0x413b2b64fD25aB1B77d117005285dd4d40c509F7")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, isContract, false, "user address should not be contract")

}
