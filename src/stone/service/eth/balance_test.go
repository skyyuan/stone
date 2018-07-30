package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErc20Balance(t *testing.T) {
	balance, err := Erc20Balance("0x413b2b64fD25aB1B77d117005285dd4d40c509F7",
		"0x54fB2B073926f20AA46604b00CcAB89E50b5505f")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, balance.Amount,
		"1583",
		"balance should be same")

	balance, err = Erc20Balance("0x413b2b64fD25aB1B77d117005285dd4d40c509F8",
		"0x54fb2b073926f20aa46604b00ccab89e50b5505f")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, balance.Amount,
		"0",
		"balance of non-existing account should be 0")

	balance, err = Erc20Balance("0x413b2b64fD25aB1B77d117005285dd4d40c509F7",
		"0x54fb2b073926f20aa46604b00ccab89e50b5505e")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, balance.Amount,
		"0",
		"balance on non-existing contract should be 0")

}

func TestEthereumBalance(t *testing.T) {
	balance, err := EthereumBalance("0xa421bb7016bf95a5c0bcdb2888c9e2b276a3403e")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, balance.Amount, "0", "balance of non-existing account should be 0")

	balance, err = EthereumBalance("0x414838e556e1623d36718a91060775c3e71c633f")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, balance.Amount, "13630000000000000000", "balance should be same")
}
