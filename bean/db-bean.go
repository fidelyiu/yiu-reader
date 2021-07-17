package bean

import (
	yiuLog "github.com/fidelyiu/yiu-go-tool/log"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB

func SetDbBean(tempDb *bbolt.DB) {
	if db == nil {
		db = tempDb
	} else {
		yiuLog.WarningLn("db-bean已经初始化!")
	}
}

func GetDbBean() *bbolt.DB {
	return db
}
