package EditSoftService

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
	EditSoftDao "yiu/yiu-reader/dao/edit-soft-dao"
	"yiu/yiu-reader/model/dto"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/response"
)

const serviceName = "编辑软件"

func Add(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var addEntity entity.EditSoft
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

	err = EditSoftDao.Save(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，数据库层错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Search(c *gin.Context) response.YiuReaderResponse {
	var searchDto dto.EditSoftSearchDto
	_ = c.ShouldBindQuery(&searchDto)
	result := response.YiuReaderResponse{}
	list, err := EditSoftDao.SearchByDto(searchDto)
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
	viewEntity, err := EditSoftDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	_ = viewEntity.CheckPath()
	result.Result = viewEntity
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Update(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var updateEntity entity.EditSoft
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
	err = EditSoftDao.Update(&updateEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，数据库层出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Delete(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	err := EditSoftDao.DeleteById(id)
	if err != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func ChangeSort(c *gin.Context, changeType enum.ChangeSortType) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	err := EditSoftDao.ChangeSort(id, changeType)
	if err != nil {
		bean.GetLoggerBean().Error("设置"+serviceName+"序号出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}
