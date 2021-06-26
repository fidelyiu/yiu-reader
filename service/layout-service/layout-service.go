package LayoutService

import (
	YiuStr "github.com/fidelyiu/yiu-go/string"
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
	maxX := YiuStr.ToInt(c.DefaultQuery("maxX", "1080"))
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
	LayoutUtil.GetDefaultLayout(currentList, &addEntity, maxX)

	err = LayoutDao.Save(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，数据库层错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Delete(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	err := LayoutDao.DeleteById(id)
	if err != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func ResizePosition(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	maxX := YiuStr.ToInt(c.DefaultQuery("maxX", "1080"))
	var resizeEntity entity.Layout
	err := c.ShouldBindJSON(&resizeEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	if resizeEntity.Id == "" {
		result.ToError("布局id不能为空")
		return result
	}
	if resizeEntity.CheckWidth() != nil || resizeEntity.CheckHeight() != nil {
		result.ToError("布局宽高无效")
		return result
	}
	// 从数据库中查找数据出来
	dbEntity, err := LayoutDao.FindById(resizeEntity.Id)
	if err != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	dbEntity.Left = resizeEntity.Left
	dbEntity.Top = resizeEntity.Top
	dbEntity.Width = resizeEntity.Width
	dbEntity.Height = resizeEntity.Height

	// 修改数据
	err = LayoutDao.Update(&dbEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	// 重置所有布局
	err = LayoutDao.FormatAll(maxX)
	if err != nil {
		bean.GetLoggerBean().Error("格式化"+serviceName+"位置出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	result.SetType(enum.ResultTypeSuccess)
	return result
}
