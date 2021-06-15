package MainController

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexHTML(c *gin.Context) {
	// c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/dist/")
	c.HTML(http.StatusOK, "index.html", nil)
}
