package NoteController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	NoteService "yiu/yiu-reader/service/note-service"
)

func Refresh(c *gin.Context) {
	NoteService.Refresh(c)
}

func SearchTree(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.SearchTree(c))
}
