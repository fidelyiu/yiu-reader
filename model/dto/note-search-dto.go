package dto

import "yiu/yiu-reader/model/enum"

type NoteSearchDto struct {
	ObjStatus   enum.ObjStatus `json:"objStatus" form:"objStatus"`     // note状态
	Path        string         `json:"path" form:"path"`               // 路径
	ParentId    string         `json:"parentId" form:"parentId"`       // 父 note ID
	Show        bool           `json:"show" form:"show"`               // show为true的note
	NotShow     bool           `json:"notShow" form:"notShow"`         // show为false的note
	WorkspaceId string         `json:"workspaceId" form:"workspaceId"` // 工作空间ID
	BadFileEnd  bool           `json:"badFileEnd" from:"badFileEnd"`   // 坏文件置后，否则只按找排序数排序，[show，已排序]>>[show，未排序]>>[notShow，已排序]>>[notShow，未排序]
}
