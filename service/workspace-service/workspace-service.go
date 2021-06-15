package WorkspaceService

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
	WorkspaceDao "yiu/yiu-reader/dao/workspace-dao"
	"yiu/yiu-reader/model/dto"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/response"
)

const serviceName = "工作空间"

func Add(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var addEntity entity.Workspace
	err := c.ShouldBindJSON(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	// 检查
	err = addEntity.Check()
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，参数检查错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	err = WorkspaceDao.Save(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，数据库层错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Search(c *gin.Context) response.YiuReaderResponse {
	var searchDto dto.WorkspaceSearchDto
	_ = c.ShouldBindQuery(&searchDto)
	result := response.YiuReaderResponse{}
	list, err := WorkspaceDao.SearchByDto(searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("获取所有"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = list
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func View(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	viewEntity, err := WorkspaceDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	_ = viewEntity.CheckPath()
	result.Result = viewEntity
	return result
}

func Update(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var updateEntity entity.Workspace
	err := c.ShouldBindJSON(&updateEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	// 检查
	err = updateEntity.Check()
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，参数检查错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = WorkspaceDao.Update(&updateEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，数据库层出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	return result
}
