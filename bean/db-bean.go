package bean

import (
	YiuLogger "github.com/fidelyiu/yiu-go/logger"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB

func SetDbBean(tempDb *bbolt.DB) {
	if db == nil {
		db = tempDb
	} else {
		YiuLogger.LogWarningLn("db-bean已经初始化!")
	}
}

func GetDbBean() *bbolt.DB {
	return db
}
