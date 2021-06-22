package LayoutService

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
	LayoutDao "yiu/yiu-reader/dao/layout-dao"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/response"
	LayoutUtil "yiu/yiu-reader/util/layout-util"
)

const serviceName = "布局"

func GetAllBySort() response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	list, err := LayoutDao.FindAllSortByUpdateTime()
	if err != nil {
		bean.GetLoggerBean().Error("获取所有"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = list
	return result
}

func Add(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var addEntity entity.Layout
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

	// 设置有效Top、Left
	currentList, err := LayoutDao.FindAllSortByUpdateTime()
	if err != nil {
		bean.GetLoggerBean().Error("获取"+serviceName+"排序列表出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	LayoutUtil.GetDefaultLayout(currentList, &addEntity)

	err = LayoutDao.Save(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，数据库层错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}
