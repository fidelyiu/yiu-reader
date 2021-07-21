package entity

import (
	"errors"
	yiuSErr "github.com/fidelyiu/yiu-go-tool/error_s"
	yiuStr "github.com/fidelyiu/yiu-go-tool/string"
	"yiu/yiu-reader/model/enum"
	PathUtil "yiu/yiu-reader/util/path-util"
)

type EditSoft struct {
	Id      string         `json:"id"`      // Uuid
	Name    string         `json:"name"`    // 名称
	Path    string         `json:"path"`    // 绝对路径
	Img     string         `json:"img"`     // 软件图标地址
	SortNum int            `json:"sortNum"` // 排序数
	Status  enum.ObjStatus `json:"status"`  // 状态
}

func (e *EditSoft) CheckPath() error {
	if !PathUtil.IsValidFile(e.Path) {
		e.Status = enum.ObjStatusInvalid
		return errors.New("编辑软件 '" + e.Path + "' 不是有效文件的绝对路径")
	}
	e.Status = enum.ObjStatusValid
	return nil
}

func (e *EditSoft) CheckName() error {
	if yiuStr.IsBlank(e.Name) {
		return errors.New("编辑软件名称不能为空")
	}
	return nil
}

func (e *EditSoft) Check() error {
	return yiuSErr.ToErrorBySep(" & ",
		e.CheckPath(),
		e.CheckName(),
	)
}
