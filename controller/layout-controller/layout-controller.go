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
