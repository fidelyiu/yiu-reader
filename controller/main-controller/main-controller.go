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

func GetOsPathSeparator(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetOsPathSeparator())
}

func GetNoteTextDocument(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetNoteTextDocument())
}

func SetNoteTextDocument(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetNoteTextDocument(c))
}

func GetNoteTextMainPoint(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetNoteTextMainPoint())
}

func SetNoteTextMainPoint(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetNoteTextMainPoint(c))
}

func GetNoteTextDir(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.GetNoteTextDir())
}

func SetNoteTextDir(c *gin.Context) {
	c.JSON(http.StatusOK, MainService.SetNoteTextDir(c))
}
