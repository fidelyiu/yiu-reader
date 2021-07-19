package NoteController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yiu/yiu-reader/model/enum"
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

func Position(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.Position(c))
}

func ChangeShow(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.ChangeShow(c))
}

func Up(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.ChangeSort(c, enum.ChangeSortTypeUp))
}

func Down(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.ChangeSort(c, enum.ChangeSortTypeDown))
}
