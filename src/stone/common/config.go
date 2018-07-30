package common

// CommConfig Configuration
type CommConfig struct {
	CmDebug        bool
	CmMysqlURL     string
	CmMysqlIdle    int
	CmMysqlMaxOpen int
	CmReqlog       string
}

func (conf *CommConfig) Debug() bool {
	return conf.CmDebug
}

func (conf *CommConfig) MysqlURL() string {
	return conf.CmMysqlURL
}

func (conf *CommConfig) MysqlIdle() int {
	return conf.CmMysqlIdle
}

func (conf *CommConfig) MysqlMaxOpen() int {
	return conf.CmMysqlMaxOpen
}

func (conf *CommConfig) Reqlog() string {
	return conf.CmReqlog
}

// Config Configuration interface
type Config interface {
	Debug() bool
	MysqlURL() string
	MysqlIdle() int
	MysqlMaxOpen() int
	Reqlog() string
}
