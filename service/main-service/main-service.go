package MainService

import (
	"errors"
	"github.com/gin-gonic/gin"
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

func SetCurrentWorkspace(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	workspaceId := c.PostForm("id")
	if workspaceId == "" {
		emptyError := errors.New("workspacePath字段不能为空")
		bean.GetLoggerBean().Error("设置当前工作空间失败!", zap.Error(emptyError))
		result.Message = emptyError.Error()
		return result
	}
	_, err := WorkspaceDao.FindById(workspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前工作空间失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetCurrentWorkspaceId(workspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前工作空间失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	return result
}
