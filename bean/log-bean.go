package bean

import (
	YiuLogger "github.com/fidelyiu/yiu-go/logger"
	"go.uber.org/zap"
)

var logger *zap.Logger

func SetLoggerBean(tempLogger *zap.Logger) {
	if logger == nil {
		logger = tempLogger
	} else {
		YiuLogger.LogWarningLn("db-bean已经初始化!")
	}
}

func GetLoggerBean() *zap.Logger {
	return logger
}
