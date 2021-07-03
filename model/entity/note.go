package entity

import (
	"yiu/yiu-reader/model/enum"
)

type Note struct {
	Id            string                // Uuid
	AbsPath       string                // 绝对路径
	Path          string                // 相对于工作空间相对路径
	Name          string                // 名称
	Alias         string                // 别名
	SortNum       int                   // 排序数
	Status        enum.ObjStatus        // 状态
	DefStatus     enum.DefinitionStatus // 是否定义过顺序，如果没定义过顺序，那么就是本地刚导入的
	WorkspaceId   string                // 所属工作空间Id
	ParentId      string                // 父级目录Id
	ParentAbsPath string                // 父级的绝对路径
	ParentPath    string                // 相对于父级的相对路径
	Level         int                   // 等级
	Show          bool                  // 是否展示
	IsDir         bool                  // 是否是文件夹
}
