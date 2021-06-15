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
