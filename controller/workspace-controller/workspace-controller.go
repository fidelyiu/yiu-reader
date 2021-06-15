package WorkspaceController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yiu/yiu-reader/model/enum"
	WorkspaceService "yiu/yiu-reader/service/workspace-service"
)

func Add(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.Add(c))
}

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.Search(c))
}

func View(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.View(c))
}

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.Update(c))
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.Delete(c))
}

func Up(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.ChangeSort(c, enum.ChangeSortTypeUp))
}

func Down(c *gin.Context) {
	c.JSON(http.StatusOK, WorkspaceService.ChangeSort(c, enum.ChangeSortTypeDown))
}
