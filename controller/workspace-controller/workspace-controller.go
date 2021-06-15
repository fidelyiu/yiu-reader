package WorkspaceController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	WorkspaceService "yiu/yiu-reader/service/workspace-service"
)

func Add(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.Add(c))
}

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.Search(c))
}
