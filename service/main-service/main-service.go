package MainService

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
	MainDao "yiu/yiu-reader/dao/main-dao"
	WorkspaceDao "yiu/yiu-reader/dao/workspace-dao"
	"yiu/yiu-reader/model/enum"
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
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetCurrentWorkspace(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	workspaceId := c.Param("id")
	if workspaceId == "" {
		emptyError := errors.New("id字段不能为空")
		bean.GetLoggerBean().Error("设置当前工作空间失败!", zap.Error(emptyError))
		result.ToError(emptyError.Error())
		return result
	}
	currentWorkspace, err := WorkspaceDao.FindById(workspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前工作空间失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = currentWorkspace.CheckPath()
	if err != nil {
		bean.GetLoggerBean().Error(workspaceId+"对应的工作空间路径无效!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetCurrentWorkspaceId(workspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前工作空间失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = currentWorkspace
	result.SetType(enum.ResultTypeSuccess)
	return result
}
