package entity

import (
	"errors"
	"yiu/yiu-reader/model/enum"
	PathUtil "yiu/yiu-reader/util/path-util"
)

type Note struct {
	Id            string         `json:"id"`            // Uuid
	AbsPath       string         `json:"absPath"`       // 绝对路径
	Path          string         `json:"path"`          // 相对于工作空间相对路径
	Name          string         `json:"name"`          // 名称
	Alias         string         `json:"alias"`         // 别名
	SortNum       int            `json:"sortNum"`       // 排序数
	Status        enum.ObjStatus `json:"status"`        // 状态
	WorkspaceId   string         `json:"workspaceId"`   // 所属工作空间Id
	ParentId      string         `json:"parentId"`      // 父级目录Id
	ParentAbsPath string         `json:"parentAbsPath"` // 父级的绝对路径
	ParentPath    string         `json:"parentPath"`    // 相对于父级的相对路径
	Level         int            `json:"level"`         // 等级
	Show          bool           `json:"show"`          // 是否展示
	IsDir         bool           `json:"isDir"`         // 是否是文件夹
	// ShowNum       int                   `json:"showNum"`       // 排除隐藏文件后的标题顺序
	// DefStatus     enum.DefinitionStatus `json:"defStatus"`     // 是否定义过顺序，如果没定义过顺序，那么就是本地刚导入的
}

func (n *Note) CheckPath() error {
	if n.IsDir {
		if !PathUtil.IsValidDir(n.AbsPath) {
			n.Status = enum.ObjStatusInvalid
			return errors.New("工作空间 '" + n.AbsPath + "' 不是有效文件夹的绝对路径")
		}
	} else {
		if !PathUtil.IsValidFile(n.AbsPath) {
			n.Status = enum.ObjStatusInvalid
			return errors.New("工作空间 '" + n.AbsPath + "' 不是有效文件的绝对路径")
		}
	}
	n.Status = enum.ObjStatusValid
	return nil
}
