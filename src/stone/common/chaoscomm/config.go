package chaoscomm

import (
	"stone/common"
	"stone/common/ethcomm"
	"stone/nsqs"
)

// ChaosConfig Configuration
type ChaosConfig struct {
	common.CommConfig
	ethcomm.EthComConfig
	Nsqconfig *nsqs.SimpleConfig
}
