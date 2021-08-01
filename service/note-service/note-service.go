package NoteService

import (
	"encoding/json"
	yiuDir "github.com/fidelyiu/yiu-go-tool/dir"
	yiuFile "github.com/fidelyiu/yiu-go-tool/file"
	yiuOs "github.com/fidelyiu/yiu-go-tool/os"
	yiuStr "github.com/fidelyiu/yiu-go-tool/string"
	yiuStrList "github.com/fidelyiu/yiu-go-tool/string_list"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"yiu/yiu-reader/bean"
	EditSoftDao "yiu/yiu-reader/dao/edit-soft-dao"
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

const serviceName = "笔记"

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
	yiuStr.OpFormatPathSeparator(&tempPath)
	if path != "" {
		tempPath += string(os.PathSeparator) + path
		yiuStr.OpFormatPathSeparator(&tempPath)
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
		if allNote[i].Level != 1 {
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

func ChangeShow(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	target, err := NoteDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("根据ID获取"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	target.Show = !target.Show
	target.SortNum = 0
	err = NoteDao.Update(&target)
	if err != nil {
		bean.GetLoggerBean().Error("更新"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Position(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	target, err := NoteDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("根据ID获取"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = target.CheckPath()
	if err != nil {
		result.ToError(err.Error())
		return result
	}

	// target.AbsPath
	err = yiuOs.DoOpenFileManagerByParent(target.AbsPath)
	if err != nil {
		bean.GetLoggerBean().Error("打开文件管理器失败!路径为\""+target.AbsPath+"\"", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func EditMarkdown(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	target, err := NoteDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("根据ID获取"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	err = target.CheckPath()
	if err != nil {
		bean.GetLoggerBean().Error(target.Name+"文件不存在!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	softId, err := MainDao.GetEditSoftId()
	if err != nil {
		bean.GetLoggerBean().Error("根据当前编辑软件ID出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	soft, err := EditSoftDao.FindById(softId)
	if err != nil {
		bean.GetLoggerBean().Error("根据当前编辑软件出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	err = soft.CheckPath()
	if err != nil {
		bean.GetLoggerBean().Error(soft.Name+"文件不存在!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	cmd := yiuOs.GetCmdWithPrefix(soft.Path + " " + target.AbsPath)
	err = cmd.Start()
	if err != nil {
		bean.GetLoggerBean().Error("使用默认编辑器启动文件异常", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func DeleteFile(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	target, err := NoteDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("根据ID获取"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	allData, err := NoteDao.FindAll()
	if err != nil {
		bean.GetLoggerBean().Error("获取所有"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	child := NoteUtil.GetChild(target, allData, false)
	err = deleteFileByTargetAndItChild(target, child)
	if err != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"文件过程中出错，稍后重试!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = NoteDao.DeleteDeepById(id)
	if err != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"记录过程中出错，稍后重试!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func DeleteBad(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")

	searchDto := dto.NoteSearchDto{}
	if id == "" {
		wId, err := MainDao.GetCurrentWorkspaceId()
		if err != nil {
			bean.GetLoggerBean().Error("获取当前工作空间失败!", zap.Error(err))
			result.ToError(err.Error())
			return result
		}
		id = wId
	}
	searchDto.WorkspaceId = id
	searchDto.ObjStatus = enum.ObjStatusInvalid
	allNote, err := NoteDao.FindBySearchDto(searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("查询所有"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	var ids []string
	for i := range allNote {
		ids = append(ids, allNote[i].Id)
	}
	delErr := NoteDao.DeleteByIds(ids)
	if delErr != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func deleteFileByTargetAndItChild(target entity.Note, child []vo.NoteTreeVo) error {
	if len(child) != 0 {
		for i := range child {
			err := deleteFileByTargetAndItChild(child[i].Data, child[i].Child)
			if err != nil {
				return err
			}
		}
	}
	if target.Status == enum.ObjStatusValid {
		// 删除文件
		err := os.Remove(target.AbsPath)
		return err
	}
	return nil
}

func Delete(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	err := NoteDao.DeleteDeepById(id)
	if err != nil {
		bean.GetLoggerBean().Error("删除"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Search(c *gin.Context) response.YiuReaderResponse {
	var searchDto dto.NoteSearchDto
	_ = c.ShouldBindQuery(&searchDto)
	result := response.YiuReaderResponse{}
	allNote, err := NoteDao.FindBySearchDto(searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("根据工作空间ID获取所有笔记失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = allNote
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func SearchTree(c *gin.Context) response.YiuReaderResponse {
	var searchDto dto.NoteSearchDto
	_ = c.ShouldBindJSON(&searchDto)

	result := response.YiuReaderResponse{}

	// 获取当前工作空间ID
	currentWorkspaceId, err := MainDao.GetCurrentWorkspaceId()
	if err != nil {
		bean.GetLoggerBean().Error("当前工作空间ID获取失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	// 获取当前工作空间
	currentWorkspace, err := WorkspaceDao.FindById(currentWorkspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("当前工作空间获取失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	if searchDto.Path != "" {
		// 找对应文件夹的笔记
		searchDto.Path = currentWorkspace.Path + string(os.PathSeparator) + searchDto.Path
		yiuStr.OpFormatPathSeparator(&searchDto.Path)
		dbNote, err := NoteDao.FindByAbsPath(searchDto.Path)
		if err != nil {
			bean.GetLoggerBean().Error("获取当前笔记ID失败!", zap.Error(err))
			result.ToError(err.Error())
			return result
		}
		searchDto.ParentId = dbNote.Id
	} else {
		// 当前工作空间所有笔记
		searchDto.WorkspaceId = currentWorkspaceId
	}
	allNote, err := NoteDao.FindBySearchDto(searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("根据工作空间ID获取所有笔记失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	if searchDto.ParentId != "" {
		searchDto.ParentId = ""
		tAllNote, allErr := NoteDao.FindBySearchDto(searchDto)
		if allErr != nil {
			bean.GetLoggerBean().Error("获取所以偶工作空间笔记失败!", zap.Error(err))
			result.ToError(allErr.Error())
			return result
		}
		tempList := NoteUtil.GetTree(allNote, searchDto.BadFileEnd)

		for i := range tempList {
			if tempList[i].Data.IsDir {
				tempList[i].Child = NoteUtil.GetChild(tempList[i].Data, tAllNote, searchDto.BadFileEnd)
			}
		}
		result.Result = tempList
	} else {
		result.Result = NoteUtil.GetTree(allNote, searchDto.BadFileEnd)
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func ChangeSort(c *gin.Context, changeType enum.ChangeSortType) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	err := NoteDao.ChangeSort(id, changeType)
	if err != nil {
		bean.GetLoggerBean().Error("设置"+serviceName+"序号出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Add(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var addEntity entity.Note
	err := c.ShouldBindJSON(&addEntity)
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	if !addEntity.IsDir {
		addEntity.Name += ".md"
	}

	// 设置 workspaceId
	if yiuStr.IsBlank(addEntity.WorkspaceId) {
		addEntity.WorkspaceId, err = MainDao.GetCurrentWorkspaceId()
		if err != nil {
			bean.GetLoggerBean().Error("获取当前工作空间ID失败!", zap.Error(err))
			result.ToError(err.Error())
			return result
		}
	}
	workspace, wErr := WorkspaceDao.FindById(addEntity.WorkspaceId)
	if wErr != nil {
		bean.GetLoggerBean().Error("无效的工作空间ID!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	// 设置 parentId
	if yiuStr.IsNotBlank(addEntity.ParentId) {
		parentEntity, pErr := NoteDao.FindById(addEntity.ParentId)
		if pErr != nil {
			bean.GetLoggerBean().Error("无效的父笔记ID!", zap.Error(err))
			result.ToError(err.Error())
			return result
		}
		addEntity.ParentAbsPath = parentEntity.AbsPath
		addEntity.AbsPath = parentEntity.AbsPath + string(os.PathSeparator) + addEntity.Name
		addEntity.Path = parentEntity.Path + string(os.PathSeparator) + addEntity.Name
		addEntity.Level = parentEntity.Level + 1
	} else {
		addEntity.ParentId = ""
		addEntity.ParentAbsPath = ""
		addEntity.AbsPath = workspace.Path + string(os.PathSeparator) + addEntity.Name
		addEntity.Path = string(os.PathSeparator) + addEntity.Name
		addEntity.Level = 1
	}
	addEntity.ParentPath = string(os.PathSeparator) + addEntity.Name
	addEntity.Show = true

	// 写入硬盘
	if addEntity.IsDir {
		if yiuDir.IsExists(addEntity.AbsPath) {
			bean.GetLoggerBean().Error(addEntity.AbsPath+"目录已存在!", zap.Error(err))
			result.ToError(addEntity.AbsPath + "目录已存在!")
			return result
		}
		dErr := yiuDir.DoMkDir(addEntity.AbsPath)
		if dErr != nil {
			bean.GetLoggerBean().Error("创建"+addEntity.AbsPath+"目录出错!", zap.Error(err))
			result.ToError(err.Error())
			return result
		}
	} else {
		if yiuFile.IsExists(addEntity.AbsPath) {
			bean.GetLoggerBean().Error(addEntity.AbsPath+"文件已存在!", zap.Error(err))
			result.ToError(addEntity.AbsPath + "文件已存在!")
			return result
		}
		fErr := yiuFile.DoCreate(addEntity.AbsPath)
		if fErr != nil {
			bean.GetLoggerBean().Error("创建"+addEntity.AbsPath+"文件出错!", zap.Error(err))
			result.ToError(err.Error())
			return result
		}
	}

	// 检查
	err = addEntity.Check()
	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，参数检查错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	err = NoteDao.Save(&addEntity)

	if err != nil {
		bean.GetLoggerBean().Error("添加"+serviceName+"出错，数据库层错误!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Update(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	var updateEntity entity.Note
	err := c.ShouldBindJSON(&updateEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错，Body参数转换出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = NoteDao.RenameNote(updateEntity)
	if err != nil {
		bean.GetLoggerBean().Error("修改"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func View(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	viewEntity, err := NoteDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = viewEntity
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func Reade(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	notePageVo := vo.NotePageVo{}
	readeEntity, err := NoteDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = readeEntity.Check()
	if err != nil {
		bean.GetLoggerBean().Error("检查"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	// 笔记信息
	notePageVo.Note = readeEntity

	file, err := os.Open(readeEntity.AbsPath)
	if err != nil {
		bean.GetLoggerBean().Error("打开"+serviceName+"文件出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	content, err := ioutil.ReadAll(file)
	if err != nil {
		bean.GetLoggerBean().Error("读取"+serviceName+"文件出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	fileInfo, err := file.Stat()
	if err != nil {
		bean.GetLoggerBean().Error("读取"+serviceName+"文件信息出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	notePageVo.ModTime = fileInfo.ModTime()
	notePageVo.Size = fileInfo.Size()

	// 内容
	notePageVo.Content = string(content)

	// 所属工作空间
	workspace, err := WorkspaceDao.FindById(readeEntity.WorkspaceId)
	if err != nil {
		bean.GetLoggerBean().Error(serviceName+"所属工作空间ID异常!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = workspace.Check()
	if err != nil {
		bean.GetLoggerBean().Error(serviceName+"所属工作空间异常!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	notePageVo.WorkSpace = workspace

	// 父路径
	notePageVo.ParentName, err = getNoteParentNamePath(readeEntity, []string{}, []string{})
	if workspace.Alias != "" {
		notePageVo.ParentName = append(notePageVo.ParentName, workspace.Alias)
	} else {
		notePageVo.ParentName = append(notePageVo.ParentName, workspace.Name)
	}

	yiuStrList.OpReverse(&notePageVo.ParentName)

	if err != nil {
		bean.GetLoggerBean().Error("读取"+serviceName+"父级名称出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	result.Result = notePageVo
	result.SetType(enum.ResultTypeSuccess)
	return result
}

func getNoteParentNamePath(entity entity.Note, parentName []string, pushId []string) ([]string, error) {
	// id避免死循环
	for i := range pushId {
		if pushId[i] == entity.Id {
			return parentName, nil
		}
	}
	pushId = append(pushId, entity.Id)
	if entity.Alias != "" {
		parentName = append(parentName, entity.Alias)
	} else {
		parentName = append(parentName, entity.Name)
	}
	if entity.ParentId != "" {
		parent, err := NoteDao.FindById(entity.ParentId)
		if err != nil {
			return nil, err
		}
		return getNoteParentNamePath(parent, parentName, pushId)
	} else {
		return parentName, nil
	}
}

func DirTree(c *gin.Context) response.YiuReaderResponse {
	result := response.YiuReaderResponse{}
	id := c.Param("id")
	readeEntity, err := NoteDao.FindById(id)
	if err != nil {
		bean.GetLoggerBean().Error("查询"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	err = readeEntity.Check()
	if err != nil {
		bean.GetLoggerBean().Error("检查"+serviceName+"出错!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	workspace, err := WorkspaceDao.FindById(readeEntity.WorkspaceId)
	if err != nil {
		bean.GetLoggerBean().Error("MD所属工作空间ID获取失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}

	var searchDto dto.NoteSearchDto
	searchDto.WorkspaceId = workspace.Id
	searchDto.Show = true
	allNote, err := NoteDao.FindBySearchDto(searchDto)
	if err != nil {
		bean.GetLoggerBean().Error("根据工作空间ID获取所有笔记失败!", zap.Error(err))
		result.ToError(err.Error())
		return result
	}
	result.Result = NoteUtil.GetTree(allNote, true)
	result.SetType(enum.ResultTypeSuccess)
	return result
}
