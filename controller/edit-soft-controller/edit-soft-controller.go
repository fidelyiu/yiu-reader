package EditSoftController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yiu/yiu-reader/model/enum"
	EditSoftService "yiu/yiu-reader/service/edit-soft-service"
)

func Add(c *gin.Context) {
	c.JSON(http.StatusOK, EditSoftService.Add(c))
}

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, EditSoftService.Search(c))
}

func View(c *gin.Context) {
	c.JSON(http.StatusOK, EditSoftService.View(c))
}

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, EditSoftService.Update(c))
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, EditSoftService.Delete(c))
}

func Up(c *gin.Context) {
	c.JSON(http.StatusOK, EditSoftService.ChangeSort(c, enum.ChangeSortTypeUp))
}

func Down(c *gin.Context) {
	c.JSON(http.StatusOK, EditSoftService.ChangeSort(c, enum.ChangeSortTypeDown))
}
