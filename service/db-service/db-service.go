package DbService

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yiu/yiu-reader/bean"
	DbDao "yiu/yiu-reader/dao/db-dao"
	"yiu/yiu-reader/model/dto"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/response"
)

const serviceName = "数据库"

func Search(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var searchDto dto.DbSearchDto
	err := c.ShouldBindJSON(&searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	searchDto.Check()

	dbSearchVo, err := DbDao.FindBySearchDto(searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	result.Result = dbSearchVo
	result.SetType(enum.ResultTypeSuccess)
	return result
}
