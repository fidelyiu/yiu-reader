package MainService

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"yiu/yiu-reader/bean"
	EditSoftDao "yiu/yiu-reader/dao/edit-soft-dao"
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
	showText, _ := MainDao.GetMainBoxShowText()
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
	showIcon, _ := MainDao.GetMainBoxShowIcon()
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
	showNum, _ := MainDao.GetMainBoxShowNum()
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
	menuOpen, _ := MainDao.GetSidebarStatus()
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

func GetEditSoft() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	editSoftId, err := MainDao.GetEditSoftId()
	if err != nil {
		bean.GetLoggerBean().Error("获取当前编辑软件失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	editSoft, err := EditSoftDao.FindById(editSoftId)
	if err != nil {
		bean.GetLoggerBean().Error("获取当前编辑软件失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	_ = editSoft.CheckPath()
	result.Result = editSoft
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetEditSoft(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	editSoftId := c.Param("id")
	if editSoftId == "" {
		emptyError := errors.New("id字段不能为空")
		bean.GetLoggerBean().Error("设置当前编辑软件失败!", zap.Error(emptyError))
		result.ToError(emptyError.Error())
		return result
	}
	currentEditSoft, err := EditSoftDao.FindById(editSoftId)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前编辑软件失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = currentEditSoft.CheckPath()
	if err != nil {
		bean.GetLoggerBean().Error(editSoftId+"对应的编辑软件路径无效!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetEditSoftId(editSoftId)
	if err != nil {
		bean.GetLoggerBean().Error("设置当前编辑软件失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = currentEditSoft
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func GetOsPathSeparator() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	result.Result = os.PathSeparator
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func GetNoteTextDocument() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	showTextDocument, err := MainDao.GetNoteTextDocument()
	if err != nil {
		result.Result = true
	} else {
		result.Result = showTextDocument
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetNoteTextDocument(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	type bodyStruct struct {
		NoteTextDocument bool `json:"noteTextDocument" form:"noteTextDocument"`
	}
	var body bodyStruct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bean.GetLoggerBean().Error("结构体绑定错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetNoteTextDocument(body.NoteTextDocument)
	if err != nil {
		bean.GetLoggerBean().Error("设置笔记页面文档工具文字是否提示错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = body.NoteTextDocument
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func GetNoteTextMainPoint() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	showTextMainPoint, err := MainDao.GetNoteTextMainPoint()
	if err != nil {
		result.Result = true
	} else {
		result.Result = showTextMainPoint
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetNoteTextMainPoint(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	type bodyStruct struct {
		NoteTextMainPoint bool `json:"noteTextMainPoint" form:"noteTextMainPoint"`
	}
	var body bodyStruct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bean.GetLoggerBean().Error("结构体绑定错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetNoteTextMainPoint(body.NoteTextMainPoint)
	if err != nil {
		bean.GetLoggerBean().Error("设置笔记页面大纲工具文字是否提示错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = body.NoteTextMainPoint
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func GetNoteTextDir() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	showTextDir, err := MainDao.GetNoteTextDir()
	if err != nil {
		result.Result = true
	} else {
		result.Result = showTextDir
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SetNoteTextDir(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	type bodyStruct struct {
		NoteTextDir bool `json:"noteTextDir" form:"noteTextDir"`
	}
	var body bodyStruct
	err := c.ShouldBindJSON(&body)
	if err != nil {
		bean.GetLoggerBean().Error("结构体绑定错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = MainDao.SetNoteTextDir(body.NoteTextDir)
	if err != nil {
		bean.GetLoggerBean().Error("设置笔记页面目录工具文字是否提示错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = body.NoteTextDir
	result.SetType(enum.ResultTypeSuccess)
	return result
}
