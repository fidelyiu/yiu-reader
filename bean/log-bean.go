package bean

import (
	yiuLog "github.com/fidelyiu/yiu-go-tool/log"
	"go.uber.org/zap"
)

var logger *zap.Logger

func SetLoggerBean(tempLogger *zap.Logger) {
	if logger == nil {
		logger = tempLogger
	} else {
		yiuLog.WarningLn("db-bean已经初始化!")
	}
}

func GetLoggerBean() *zap.Logger {
	return logger
}
