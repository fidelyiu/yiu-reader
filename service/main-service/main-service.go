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

func GetMainBoxShowText() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	showText, err := MainDao.GetMainBoxShowText()
	if err != nil {
		bean.GetLoggerBean().Error("获取当前主盒子是否展示提示错误!", zap.Error(err))
	}
	result.Result = showText
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetMainBoxShowText(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	type bodyStruct struct {
		ShowText bool `json:"showText" form:"showText"`
	}
	var body bodyStruct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bean.GetLoggerBean().Error("结构体绑定错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetMainBoxShowText(body.ShowText)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前主盒子是否展示提示错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = body.ShowText
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func GetMainBoxShowIcon() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	showIcon, err := MainDao.GetMainBoxShowIcon()
	if err != nil {
		bean.GetLoggerBean().Error("获取当前主盒子是否展示提示Icon错误!", zap.Error(err))
	}
	result.Result = showIcon
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetMainBoxShowIcon(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	type bodyStruct struct {
		ShowIcon bool `json:"showIcon" form:"showIcon"`
	}
	var body bodyStruct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bean.GetLoggerBean().Error("结构体绑定错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetMainBoxShowIcon(body.ShowIcon)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前主盒子是否展示提示Icon错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = body.ShowIcon
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func GetMainBoxShowNum() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	showNum, err := MainDao.GetMainBoxShowNum()
	if err != nil {
		bean.GetLoggerBean().Error("获取当前主盒子是否展示提示序号错误!", zap.Error(err))
	}
	result.Result = showNum
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetMainBoxShowNum(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	type bodyStruct struct {
		ShowNum bool `json:"showNum" form:"showNum"`
	}
	var body bodyStruct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bean.GetLoggerBean().Error("结构体绑定错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetMainBoxShowNum(body.ShowNum)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前主盒子是否展示提示序号错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = body.ShowNum
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func GetSidebarStatus() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	menuOpen, err := MainDao.GetSidebarStatus()
	if err != nil {
		bean.GetLoggerBean().Error("获取当前菜单是否展示错误!", zap.Error(err))
	}
	result.Result = menuOpen
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetSidebarStatus(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	type bodyStruct struct {
		SidebarStatus bool `json:"sidebarStatus" form:"sidebarStatus"`
	}
	var body bodyStruct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bean.GetLoggerBean().Error("结构体绑定错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetSidebarStatus(body.SidebarStatus)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前菜单是否展示错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = body.SidebarStatus
	result.SetType(enum.ResultTypeSuccess)
	return result
}
