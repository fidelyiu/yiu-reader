package entity

import (
	"errors"
	YiuErrorList "github.com/fidelyiu/yiu-go/error_list"
	YiuStr "github.com/fidelyiu/yiu-go/string"
	"yiu/yiu-reader/model/enum"
	PathUtil "yiu/yiu-reader/util/path-util"
)

type Workspace struct {
	Id      string         `json:"id"`      // Uuid
	Path    string         `json:"path"`    // 绝对路径
	Name    string         `json:"name"`    // 名称
	SortNum int            `json:"sortNum"` // 排序数
	Status  enum.ObjStatus `json:"status"`  // 状态
}

func (w *Workspace) CheckPath() error {
	if !PathUtil.IsValidDir(w.Path) {
		w.Status = enum.ObjStatusInvalid
		return errors.New("工作空间 '" + w.Path + "' 不是有效绝对路径")
	}
	w.Status = enum.ObjStatusValid
	return nil
}

func (w *Workspace) CheckName() error {
	if YiuStr.IsBlank(w.Name) {
		return errors.New("工作空间名称不能为空")
	}
	return nil
}

func (w *Workspace) Check() error {
	return YiuErrorList.ToError([]error{
		w.CheckPath(),
		w.CheckName(),
	}, " & ")
}
