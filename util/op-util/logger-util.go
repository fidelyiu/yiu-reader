package OpUtil

import (
	YiuLogger "github.com/fidelyiu/yiu-go/logger"
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
)

func CreateLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		YiuLogger.LogErrorLn("日志初始化出错：")
		YiuLogger.LogErrorLn(err)
		return
	}
	bean.SetLoggerBean(logger)
}
