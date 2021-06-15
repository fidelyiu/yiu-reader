package MainService

import (
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
	MainDao "yiu/yiu-reader/dao/main-dao"
	WorkspaceDao "yiu/yiu-reader/dao/workspace-dao"
	"yiu/yiu-reader/model/response"
)

func GetCurrentWorkspace() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	currentWorkspaceId, err := MainDao.GetCurrentWorkspaceId()
	if err != nil {
		bean.GetLoggerBean().Error("获取当前工作空间Path失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	workspace, err := WorkspaceDao.FindById(currentWorkspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("获取当前工作空间失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	_ = workspace.CheckPath()
	result.Result = workspace
	return result
}
