package ImageController

import (
	"github.com/gin-gonic/gin"
	ImageService "yiu/yiu-reader/service/image-service"
)

func Get(c *gin.Context) {
	ImageService.Get(c)
}

func Load(c *gin.Context) {
	ImageService.Load(c)
}
