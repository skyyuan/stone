package syncdb

import (
	"stone/common/ethcomm"
	"stone/service/eth"
)

func init() {
	config := &ethcomm.EthConfig{
		EthComGethEndpoint: eth.RopstenTestNode,
	}
	eth.ServiceInit(config)
}
