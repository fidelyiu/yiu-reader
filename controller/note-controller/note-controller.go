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

func Search(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.Search(c))
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.Delete(c))
}

func DeleteFile(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.DeleteFile(c))
}
