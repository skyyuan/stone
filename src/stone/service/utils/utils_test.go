package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaddingData(t *testing.T) {
	paddingStr := PaddingData("a9059cbb", "0x31EFd75bc0b5fbafc6015Bd50590f4fDab6a3F22",
		"0x29A2241AF62C0000")
	assert.Equal(t, "0xa9059cbb00000000000000000000000031EFd75bc0b5fbafc6015Bd50590f4fDab6a3F2200000000000000000000000000000000000000000000000029A2241AF62C0000", paddingStr, "Does not meet expectation")

	paddingStr = PaddingData("a9059cbb", "31EFd75bc0b5fbafc6015Bd50590f4fDab6a3F22",
		"29A2241AF62C0000")
	assert.Equal(t, "0xa9059cbb00000000000000000000000031EFd75bc0b5fbafc6015Bd50590f4fDab6a3F2200000000000000000000000000000000000000000000000029A2241AF62C0000", paddingStr, "Does not meet expectation")
}

func TestParseBig256(t *testing.T) {
	expBigInt, _ := new(big.Int).SetString("8233858825196223825349808356040661660160778559683678103130061641281455237550", 10)
	bigInteger, _ := ParseBig256("0x123432edfbac45436778978abcdef1b3c7d8eeff900d1a88c274deffedcba9ae")
	assert.Equal(t, expBigInt, bigInteger, "it should success to parse int that is 256bits")

	expBigInt, _ = new(big.Int).SetString("123432edfbac45436778978abcdef1b3c7d8eeff900d1a88c274deffedcba9ae", 16)
	bigInteger, _ = ParseBig256("8233858825196223825349808356040661660160778559683678103130061641281455237550")
	assert.Equal(t, expBigInt, bigInteger, "it should success to parse int that is 256bits")

	bigInteger, _ = ParseBig256("123432edfbac45436778978abcdef1b3c7d8eeff900d1a88c274deffedcba9ae")
	assert.Nil(t, bigInteger, "it should fail to parse hex without 0x prefix")

	bigInteger, _ = ParseBig256("0x123432edfbac45436778978abcdef1b3c7d8eeff900d1a88c274deffedcba9ae1")
	assert.Nil(t, bigInteger, "it should fail to parse int that is more than 256bits")
}

func TestHexToInt(t *testing.T) {
	expValue := int64(1489509580025388766)
	value, _ := HexToInt("0x14ABCD123ef45ede")
	assert.Equal(t, expValue, value, "result does not meet expectation")

	value, _ = HexToInt("14ABCD123ef45ede")
	assert.Equal(t, expValue, value, "result does not meet expectation")
}

func TestIntToHex(t *testing.T) {
	str := IntToHex(1489509580025388766)
	assert.Equal(t, "0x14abcd123ef45ede", str, "result does not meet expectation")
}

func TestInt64ToHex(t *testing.T) {
	str := Int64ToHex(1489509580025388766)
	assert.Equal(t, "0x14abcd123ef45ede", str, "result does not meet expectation")
}
