package ImageController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ImageService "yiu/yiu-reader/service/image-service"
)

func Upload(c *gin.Context) {
	c.JSON(http.StatusOK, ImageService.Upload(c))
}
func Load(c *gin.Context) {
	ImageService.Load(c)
}
