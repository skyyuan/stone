package etherscan

import (
	"testing"
)

func TestSendRawTransaction(t *testing.T) {

	hexStr := "0xf86b3b84b2d05e008252089485cafe9c371791487ef41588e6ab40dd4ec4185188016345785d8a00008029a06bfda0fb3ae65d974710eb05ec8f6947c73111b1c8ad3bb4ffb9a22ddafc685ea0022f9710158e1e9c1b71e08804e2d081b9a8a59731e23f0eab2cca2227b8e2b2"
	SendRawTransaction(true, hexStr)
}
