package vo

import (
	"time"
	"yiu/yiu-reader/model/entity"
)

// NotePageVo 阅读笔记页面的VO
type NotePageVo struct {
	Content    string           `json:"content"`    // 笔记内容
	Note       entity.Note      `json:"note"`       // 笔记实体
	ParentName []string         `json:"parentName"` // 父级路径，从工作空间开始
	WorkSpace  entity.Workspace `json:"workspace"`  // 所属工作空间
	ModTime    time.Time        `json:"modTime"`    // 最后修改时间
	Size       int64            `json:"size"`       // 大小
}
