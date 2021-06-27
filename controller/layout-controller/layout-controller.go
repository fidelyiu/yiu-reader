package LayoutController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	LayoutService "yiu/yiu-reader/service/layout-service"
)

func Add(c *gin.Context) {
	c.JSON(http.StatusOK, LayoutService.Add(c))
}

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, LayoutService.GetAllBySort())
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, LayoutService.Delete(c))
}

func ResizePosition(c *gin.Context) {
	c.JSON(http.StatusOK, LayoutService.ResizePosition(c))
}

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, LayoutService.Update(c))
}

func View(c *gin.Context) {
	c.JSON(http.StatusOK, LayoutService.View(c))
}
