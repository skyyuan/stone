package webcomm

import (
	"stone/common"
	"stone/nsqs"
)

// Config Configuration
type Config struct {
	common.CommConfig
	Nsqconfig *nsqs.SimpleConfig
}
