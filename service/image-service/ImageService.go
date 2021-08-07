package ImageService

import (
	yiuFile "github.com/fidelyiu/yiu-go-tool/file"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	ImageDao "yiu/yiu-reader/dao/image-dao"
	NoteDao "yiu/yiu-reader/dao/note-dao"
	"yiu/yiu-reader/model/enum"
)

func Get(c *gin.Context) {
	src := c.Query("src")
	noteId := c.Query("id")
	srcPath := ""
	if src != "" {
		if yiuFile.IsExists(src) {
			srcPath = src
		} else {
			if noteId != "" {
				noteEntity, err := NoteDao.FindById(noteId)
				if err != nil {
					return
				}
				srcPath = filepath.Join(noteEntity.AbsPath, src)
			}
		}
	}
	if srcPath == "" {
		return
	}
	if !yiuFile.IsExists(srcPath) {
		return
	}
	imageEntity, _ := ImageDao.FindBySrc(srcPath)
	if imageEntity.Status != enum.ObjStatusValid {
		err := ImageDao.SaveBySrc(srcPath)
		if err != nil {
			return
		}
		imageEntity, _ = ImageDao.FindBySrc(srcPath)
	}
	if imageEntity.Id == "" {
		return
	}
	c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/image/"+imageEntity.Path)
}

func Load(c *gin.Context) {
	src := c.Query("src")
	noteId := c.Query("id")
	srcPath := ""
	if src != "" {
		if yiuFile.IsExists(src) {
			srcPath = src
		} else {
			if noteId != "" {
				noteEntity, err := NoteDao.FindById(noteId)
				if err != nil {
					return
				}
				srcPath = filepath.Join(noteEntity.AbsPath, src)
			}
		}
	}
	if srcPath == "" {
		return
	}
	if !yiuFile.IsExists(srcPath) {
		return
	}
	c.File(srcPath)
	// f, err := os.Open(srcPath)
	// if err != nil {
	// 	c.Writer.WriteHeader(http.StatusNotFound)
	// 	c.handlers = group.engine.noRoute
	// 	// Reset index
	// 	c.index = -1
	// 	return
	// }
	// f.Close()
}
