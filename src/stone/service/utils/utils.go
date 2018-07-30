package utils

import (
	"math/big"
	"strconv"
	"strings"
)

const (
	// HexPrefix 16进制前缀
	HexPrefix = "0x"
)

var paddingstr = "0000000000000000000000000000000000000000000000000000000000000000" // len is 64

// PaddingData padding data for transaction input data
func PaddingData(method string, params ...string) string {
	var res string
	if !strings.HasPrefix(method, HexPrefix) {
		res = HexPrefix + method
	}
	for _, item := range params {
		if strings.HasPrefix(item, HexPrefix) {
			item = item[2:]
		}
		paddingString := paddingstr[:64-len(item)]
		tmp := string(paddingString) + item
		res += tmp
	}
	return res
}

// ParseBig256 parses s as a 256 bit integer in decimal or hexadecimal syntax.
// Leading zeros are accepted. The empty string parses as zero.
// This function is copied from go-ethereum/common/math/big.go.
func ParseBig256(s string) (*big.Int, bool) {
	if s == "" {
		return new(big.Int), true
	}
	var bigint *big.Int
	var ok bool
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		bigint, ok = new(big.Int).SetString(s[2:], 16)
	} else {
		bigint, ok = new(big.Int).SetString(s, 10)
	}
	if ok && bigint.BitLen() > 256 {
		bigint, ok = nil, false
	}
	return bigint, ok
}

// HexToInt hex string to int
func HexToInt(hexStr string) (num int64, err error) {
	tmpStr := hexStr
	if strings.HasPrefix(hexStr, "0x") {
		tmpStr = hexStr[2:]
	}
	num, err = strconv.ParseInt(tmpStr, 16, 64)
	return
}

// IntToHex int to hex
func IntToHex(num uint64) string {
	numStr := strconv.FormatUint(num, 16)
	return "0x" + numStr
}

// Int64ToHex hh
func Int64ToHex(num int64) string {
	numStr := strconv.FormatInt(num, 16)
	return "0x" + numStr
}
