package WorkspaceService

import (
	"encoding/json"
	YiuStr "github.com/fidelyiu/yiu-go/string"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"yiu/yiu-reader/bean"
	MainDao "yiu/yiu-reader/dao/main-dao"
	NoteDao "yiu/yiu-reader/dao/note-dao"
	WorkspaceDao "yiu/yiu-reader/dao/workspace-dao"
	"yiu/yiu-reader/model/dto"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/response"
	"yiu/yiu-reader/model/vo"
	NoteUtil "yiu/yiu-reader/util/note-util"
)

const serviceName = "工作空间"

func Add(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var addEntity entity.Workspace
	err := c.ShouldBindJSON(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	// 检查
	err = addEntity.Check()
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，参数检查错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	err = WorkspaceDao.Save(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，数据库层错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Search(c *gin.Context) response.YiuReaderResponse {
	var searchDto dto.WorkspaceSearchDto
	_ = c.ShouldBindQuery(&searchDto)
	result := response.YiuReaderResponse{}
	list, err := WorkspaceDao.SearchByDto(searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("获取所有"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = list
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func View(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	viewEntity, err := WorkspaceDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	_ = viewEntity.CheckPath()
	result.Result = viewEntity
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Update(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var updateEntity entity.Workspace
	err := c.ShouldBindJSON(&updateEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	// 检查
	err = updateEntity.Check()
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，参数检查错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = WorkspaceDao.Update(&updateEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，数据库层出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Delete(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	err := WorkspaceDao.DeleteById(id)
	if err != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func ChangeSort(c *gin.Context, changeType enum.ChangeSortType) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	err := WorkspaceDao.ChangeSort(id, changeType)
	if err != nil {
		bean.GetLoggerBean().Error("设置"+serviceName+"序号出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Refresh(c *gin.Context) {
	path := c.Query("path")
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)

	// _ = ws.WriteMessage(websocket.TextMessage, []byte("pong"))

	// 获取当前工作空间ID
	currentWorkspaceId, err := MainDao.GetCurrentWorkspaceId()
	if err != nil {
		bean.GetLoggerBean().Error("当前工作空间ID获取失败!", zap.Error(err))
		return
	}

	// 获取当前工作空间
	currentWorkspace, err := WorkspaceDao.FindById(currentWorkspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("当前工作空间获取失败!", zap.Error(err))
		return
	}

	// 工具当前工作空间的路径获取所有数据
	tempPath := currentWorkspace.Path
	tempLevel := 1
	YiuStr.OpFormatPathSeparator(&tempPath)
	if path != "" {
		tempPath += string(os.PathSeparator) + path
		YiuStr.OpFormatPathSeparator(&tempPath)
		dbNote, err := NoteDao.FindByAbsPath(tempPath)
		if err != nil {
			bean.GetLoggerBean().Error("获取当前笔记等级失败!", zap.Error(err))
			return
		}
		tempLevel = dbNote.Level + 1
	}

	files, err := ioutil.ReadDir(tempPath)
	if err != nil {
		bean.GetLoggerBean().Error("读取工作空间路径失败!", zap.Error(err))
		return
	}

	// 获取所有Note
	allNote, err := NoteDao.FindAll()
	if err != nil {
		bean.GetLoggerBean().Error("获取所有笔记失败!", zap.Error(err))
		return
	}

	// 开始读取的管道通知
	startChan := make(chan vo.NoteReadVo, 50)
	// 读取结束的管道通知
	endChan := make(chan vo.NoteReadVo, 50)
	stopCh := make(chan struct{})
	// 接收结果
	var result []entity.Note
	var resultLock sync.Mutex
	// 用于记录所有执行 OutNoteByPath 的 goroutine
	var noteWg sync.WaitGroup
	var stopWg sync.WaitGroup

	for _, v := range files {
		noteWg.Add(1)
		go NoteUtil.OutNoteByPath(tempPath+string(os.PathSeparator)+v.Name(),
			currentWorkspace, v,
			startChan, endChan,
			&stopWg,
			&result, &resultLock,
			&noteWg, tempLevel,
			allNote)
	}
	go func() {
		for {
			select {
			case i := <-startChan:
				buf, _ := json.Marshal(i)
				_ = ws.WriteMessage(websocket.TextMessage, buf)
			case i := <-endChan:
			priority1:
				for {
					select {
					case ti := <-startChan:
						buf, _ := json.Marshal(ti)
						_ = ws.WriteMessage(websocket.TextMessage, buf)
					default:
						break priority1
					}
				}
				buf, _ := json.Marshal(i)
				_ = ws.WriteMessage(websocket.TextMessage, buf)
				stopWg.Done()
			case <-stopCh:
				return
			}
		}
	}()
	// 确保 goroutine 都执行完成
	noteWg.Wait()
	// 确保所有 endChan 都被处理
	stopWg.Wait()
	stopCh <- struct{}{}
	// 将所有 noteResult 保存起来
	err = NoteDao.SaveAll(result)
	if err != nil {
		bean.GetLoggerBean().Error("保存所有笔记失败!", zap.Error(err))
		return
	}

	// 更新所有父ID
	allNote, err = NoteDao.FindAll()
	if err != nil {
		bean.GetLoggerBean().Error("获取所有笔记失败!", zap.Error(err))
		return
	}
	for i := range allNote {
		if allNote[i].ParentId == "" && allNote[i].Level != 1 {
			parentEntity, err := NoteDao.FindByAbsPath(allNote[i].ParentAbsPath)
			if err != nil {
				bean.GetLoggerBean().Error("根据"+allNote[i].ParentAbsPath+"找不到笔记!", zap.Error(err))
				continue
			}
			allNote[i].ParentId = parentEntity.Id
			err = NoteDao.Update(&allNote[i])
			if err != nil {
				bean.GetLoggerBean().Error("修改笔记错误!", zap.Error(err))
				continue
			}
		}
	}
}
