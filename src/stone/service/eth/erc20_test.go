package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetERC20Info(t *testing.T) {
	res, err := GetERC20Info("0x937ea0eed76dc56a1862d473f6b5dbf6ea1df191")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.Name, "OMGToken")

}

func TestConvertToAscii(t *testing.T) {
	tmp := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000084f4d47546f6b656e000000000000000000000000000000000000000000000000"
	res := convertToASCII(tmp)
	assert.Equal(t, res, "OMGToken", "should be consistent")
}
