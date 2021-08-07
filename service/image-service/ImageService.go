package ImageService

import (
	yiuDir "github.com/fidelyiu/yiu-go-tool/dir"
	yiuFile "github.com/fidelyiu/yiu-go-tool/file"
	yiuTime "github.com/fidelyiu/yiu-go-tool/time"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"path/filepath"
	"yiu/yiu-reader/bean"
	NoteDao "yiu/yiu-reader/dao/note-dao"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/response"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

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
				srcPath = filepath.Join(noteEntity.AbsPath, "..", src)
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
}

func Upload(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	// 单文件
	file, _ := c.FormFile("file")
	// filepath.Ext(file.Filename)
	file.Filename = yiuTime.GetNowStr22() + filepath.Ext(file.Filename)
	// 上传文件至指定目录
	path := FieldUtil.ImageAdd + yiuTime.GetNowStr21()
	err := yiuDir.DoMkDirAll(path)
	if err != nil {
		bean.GetLoggerBean().Error("创建"+path+"文件夹出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		bean.GetLoggerBean().Error("上传文件出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}
