package DbController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	DbService "yiu/yiu-reader/service/db-service"
)

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, DbService.Search(c))
}
