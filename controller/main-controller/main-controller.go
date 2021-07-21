package MainController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	MainService "yiu/yiu-reader/service/main-service"
)

func IndexHTML(c *gin.Context) {
	// c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/dist/")
	c.HTML(http.StatusOK, "index.html", nil)
}

func GetCurrentWorkspace(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetCurrentWorkspace())
}

func SetCurrentWorkspace(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetCurrentWorkspace(c))
}

func GetMainBoxShowText(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetMainBoxShowText())
}

func SetMainBoxShowText(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetMainBoxShowText(c))
}

func GetMainBoxShowIcon(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetMainBoxShowIcon())
}

func SetMainBoxShowIcon(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetMainBoxShowIcon(c))
}

func GetMainBoxShowNum(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetMainBoxShowNum())
}

func SetMainBoxShowNum(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetMainBoxShowNum(c))
}

func GetSidebarStatus(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetSidebarStatus())
}

func SetSidebarStatus(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetSidebarStatus(c))
}

func GetEditSoft(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetEditSoft())
}

func SetEditSoft(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetEditSoft(c))
}
