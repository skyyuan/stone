package main

import "stone/service/eth/syncdb"

func dbMigrate() {
	syncdb.MigrateDB()
}
