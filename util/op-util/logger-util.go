package OpUtil

import (
	yiuLog "github.com/fidelyiu/yiu-go-tool/log"
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
)

func CreateLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		yiuLog.ErrorLn("日志初始化出错：")
		yiuLog.ErrorLn(err)
		return
	}
	bean.SetLoggerBean(logger)
}
