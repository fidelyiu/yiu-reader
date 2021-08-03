package NoteUtil

import (
	"github.com/go-basic/uuid"
	"go.uber.org/zap"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/vo"
)

// OutNoteByPath 判断当前路径下的Note实体
// - path：当前文件路径，绝对路径
// - workspacePath：工作空间
// - file：当前文件信息
// - startChan：开始通道
// - endChan：结束通道
// - stopWg：没加入一个start就加一
// - noteWg：标识所有 OutNoteByPath 方法执行完成的 wg
// - level：当前文件的等级
// - noteList：当前所有 Note，用于判断当前result。
func OutNoteByPath(path string,
	workspace entity.Workspace,
	file fs.FileInfo,
	startChan chan<- vo.NoteReadVo, endChan chan<- vo.NoteReadVo,
	stopWg *sync.WaitGroup,
	result *[]entity.Note, resultLock *sync.Mutex,
	noteWg *sync.WaitGroup, level int,
	noteList []entity.Note) {
	defer noteWg.Done()
	if !file.IsDir() {
		if strings.HasSuffix(file.Name(), ".md") {
			startChan <- vo.NoteReadVo{
				Path:   path,
				Result: enum.NoteReadResultStart,
			}
			stopWg.Add(1)
			tempEntity := entity.Note{
				Id:            uuid.New(),
				AbsPath:       path,
				Path:          strings.TrimPrefix(path, workspace.Path),
				Name:          file.Name(),
				WorkspaceId:   workspace.Id,
				ParentPath:    string(os.PathSeparator) + file.Name(),
				ParentAbsPath: strings.TrimSuffix(strings.TrimPrefix(path, workspace.Path), string(os.PathSeparator)+file.Name()),
				Level:         level,
				Show:          true,
				IsDir:         false,
			}
			tempResult := noteListIsInclude(noteList, tempEntity)
			if tempResult == enum.NoteReadResultNotImport {
				resultLock.Lock()
				*result = append(*result, tempEntity)
				resultLock.Unlock()
			}
			endChan <- vo.NoteReadVo{
				Path:   path,
				Result: tempResult,
			}
		}
		return
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		bean.GetLoggerBean().Error("文件读取失败!", zap.Error(err))
		endChan <- vo.NoteReadVo{
			Path:   path,
			Result: enum.NoteReadResultFail,
		}
		return
	} else {
		startChan <- vo.NoteReadVo{
			Path:   path,
			Result: enum.NoteReadResultStart,
		}
		stopWg.Add(1)
		tempEntity := entity.Note{
			Id:            uuid.New(),
			AbsPath:       path,
			Path:          strings.TrimPrefix(path, workspace.Path),
			Name:          file.Name(),
			WorkspaceId:   workspace.Id,
			ParentPath:    string(os.PathSeparator) + file.Name(),
			ParentAbsPath: strings.TrimSuffix(path, string(os.PathSeparator)+file.Name()),
			Level:         level,
			Show:          true,
			IsDir:         true,
		}
		tempResult := noteListIsInclude(noteList, tempEntity)
		if tempResult == enum.NoteReadResultNotImport {
			resultLock.Lock()
			*result = append(*result, tempEntity)
			resultLock.Unlock()
		}
		endChan <- vo.NoteReadVo{
			Path:   path,
			Result: tempResult,
		}
	}
	for _, v := range files {
		if v.IsDir() {
			noteWg.Add(1)
			go OutNoteByPath(path+string(os.PathSeparator)+v.Name(),
				workspace, v,
				startChan, endChan,
				stopWg,
				result, resultLock,
				noteWg, level+1,
				noteList)
		} else {
			if strings.HasSuffix(v.Name(), ".md") {
				startChan <- vo.NoteReadVo{
					Path:   path + string(os.PathSeparator) + v.Name(),
					Result: enum.NoteReadResultStart,
				}
				stopWg.Add(1)
				tempEntity := entity.Note{
					Id:            uuid.New(),
					AbsPath:       path + string(os.PathSeparator) + v.Name(),
					Path:          strings.TrimPrefix(path+string(os.PathSeparator)+v.Name(), workspace.Path),
					Name:          v.Name(),
					WorkspaceId:   workspace.Id,
					ParentPath:    string(os.PathSeparator) + v.Name(),
					ParentAbsPath: path,
					Level:         level + 1,
					Show:          true,
					IsDir:         false,
				}
				tempResult := noteListIsInclude(noteList, tempEntity)
				if tempResult == enum.NoteReadResultNotImport {
					resultLock.Lock()
					*result = append(*result, tempEntity)
					resultLock.Unlock()
				}
				endChan <- vo.NoteReadVo{
					Path:   path + string(os.PathSeparator) + v.Name(),
					Result: tempResult,
				}
			}
		}
	}
}

// noteListIsInclude 判断当前笔记是否在noteList中
func noteListIsInclude(noteList []entity.Note, it entity.Note) enum.NoteReadResult {
	if it.AbsPath == "" {
		return enum.NoteReadResultFail
	}
	if len(noteList) == 0 {
		return enum.NoteReadResultNotImport
	}
	for _, v := range noteList {
		if v.AbsPath == it.AbsPath {
			return enum.NoteReadResultImport
		}
	}
	return enum.NoteReadResultNotImport
}

func GetTree(noteList []entity.Note, badFileEnd bool) []vo.NoteTreeVo {
	var result []vo.NoteTreeVo
	if len(noteList) == 0 {
		return result
	}
	// 先找出最高的等级
	level := noteList[0].Level
	for i := range noteList {
		if level > noteList[i].Level {
			level = noteList[i].Level
		}
	}

	for i := range noteList {
		if level == noteList[i].Level {
			result = append(result, vo.NoteTreeVo{
				Data: noteList[i],
			})
		}
	}

	// 先将最高级的文件排序
	if badFileEnd {
		// [show，已排序]>>[show，未排序]>>[notShow，已排序]>>[notShow，未排序]
		// 最后考虑排序数
		sort.Slice(result, func(i, j int) bool {
			// Show不一样
			if result[i].Data.Show != result[j].Data.Show {
				// true则i排在前面，false则i排在后面
				return result[i].Data.Show
			}
			// Show一样，是否排序都大于0
			if (result[i].Data.SortNum > 0) != (result[j].Data.SortNum > 0) {
				// 大于0则排前面，小于等于0则排后面
				return result[i].Data.SortNum > 0
			}
			return result[i].Data.SortNum < result[j].Data.SortNum
		})
	} else {
		// 只考虑排序数
		sort.Slice(result, func(i, j int) bool {
			return result[i].Data.SortNum < result[j].Data.SortNum
		})
	}

	for i := range result {
		if result[i].Data.IsDir {
			result[i].Child = GetChild(result[i].Data, noteList, badFileEnd)
		}
	}
	return result
}

func GetChild(parent entity.Note, noteList []entity.Note, badFileEnd bool) []vo.NoteTreeVo {
	var result []vo.NoteTreeVo
	if len(noteList) == 0 {
		return result
	}
	for i := range noteList {
		if noteList[i].ParentId == parent.Id {
			result = append(result, vo.NoteTreeVo{
				Data:  noteList[i],
				Child: GetChild(noteList[i], noteList, badFileEnd),
			})
		}
	}
	if len(result) != 0 {
		if badFileEnd {
			// [show，已排序]>>[show，未排序]>>[notShow，已排序]>>[notShow，未排序]
			// 最后考虑排序数
			sort.Slice(result, func(i, j int) bool {
				// Show不一样
				if result[i].Data.Show != result[j].Data.Show {
					// true则i排在前面，false则i排在后面
					return result[i].Data.Show
				}
				// Show一样，是否排序都大于0
				if (result[i].Data.SortNum > 0) != (result[j].Data.SortNum > 0) {
					// 大于0则排前面，小于等于0则排后面
					return result[i].Data.SortNum > 0
				}
				return result[i].Data.SortNum < result[j].Data.SortNum
			})
		} else {
			// 只考虑排序数
			sort.Slice(result, func(i, j int) bool {
				return result[i].Data.SortNum < result[j].Data.SortNum
			})
		}
	}
	return result
}
