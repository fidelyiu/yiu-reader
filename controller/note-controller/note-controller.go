package NoteController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yiu/yiu-reader/model/enum"
	NoteService "yiu/yiu-reader/service/note-service"
)

func Add(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.Add(c))
}

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.Update(c))
}

func View(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.View(c))
}

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

func DeleteBad(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.DeleteBad(c))
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

func EditMarkdown(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.EditMarkdown(c))
}

func Reade(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.Reade(c))
}

func DirTree(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.DirTree(c))
}

func GetNumDocument(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.GetNumDocument(c))
}

func SetNumDocument(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.SetNumDocument(c))
}

func GetNumMainPoint(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.GetNumMainPoint(c))
}

func SetNumMainPoint(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.SetNumMainPoint(c))
}

func GetNumDir(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.GetNumDir(c))
}

func SetNumDir(c *gin.Context) {
	c.JSON(http.StatusOK, NoteService.SetNumDir(c))
}
