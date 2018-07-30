package ethcomm

import (
	"stone/common"
	"stone/nsqs"
)

// EthConfig ethereum Configuration
type EthConfig struct {
	common.CommConfig
	EthComConfig
	Nsqconfig *nsqs.SimpleConfig
}

// EthComConfig config
type EthComConfig struct {
	EthComMysqlURL     string
	EthComMysqlIdle    int
	EthComMysqlMaxOpen int
	EthComGethEndpoint string
	EthComTestNet      bool
	DisableAutoErc20   bool
}

// EthMysqlURL interface implementation
func (c *EthComConfig) EthMysqlURL() string {
	return c.EthComMysqlURL
}

// EthMysqlIdle interface implementation
func (c *EthComConfig) EthMysqlIdle() int {
	return c.EthComMysqlIdle
}

// EthMysqlMaxOpen interface implementation
func (c *EthComConfig) EthMysqlMaxOpen() int {
	return c.EthComMysqlMaxOpen
}

// GethEndpoint interface implementation
func (c *EthComConfig) GethEndpoint() string {
	return c.EthComGethEndpoint
}

// EthTestNet interface implementation
func (c *EthComConfig) EthTestNet() bool {
	return c.EthComTestNet
}

// Config Configuration interface
type Config interface {
	Debug() bool
	EthMysqlURL() string
	EthMysqlIdle() int
	EthMysqlMaxOpen() int
	GethEndpoint() string
	EthTestNet() bool
}
