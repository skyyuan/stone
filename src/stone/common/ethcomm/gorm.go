package ethcomm

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// DB gorm数据库实例
var DB *gorm.DB

// GormDB 封装后的gorm数据库实例
type GormDB struct {
	*gorm.DB
	gdbDone bool
}

// InitDB 初始化数据库
func InitDB(config Config) {
	// var connstring string
	idb, err := gorm.Open("mysql", config.EthMysqlURL())
	if err != nil {
		panic(err)
	}
	// Then you could invoke `*sql.DB`'s functions with it
	idb.DB().SetMaxIdleConns(config.EthMysqlIdle())
	idb.DB().SetMaxOpenConns(config.EthMysqlMaxOpen())
	idb.LogMode(config.Debug())

	DB = idb
}

// DBClose 关闭数据库
func DBClose() {
	DB.Close()
}

// DBBegin 打开一个transaction
func DBBegin() *GormDB {
	txn := DB.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	return &GormDB{txn, false}
}

// DBCommit 提交并关闭transaction
func (c *GormDB) DBCommit() {
	if c.gdbDone {
		return
	}
	tx := c.Commit()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}

// DBRollback 回滚并关闭transaction
func (c *GormDB) DBRollback() {
	if c.gdbDone {
		return
	}
	tx := c.Rollback()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}
