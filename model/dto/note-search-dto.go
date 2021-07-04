package dto

import "yiu/yiu-reader/model/enum"

type NoteSearchDto struct {
	ObjStatus   enum.ObjStatus `json:"objStatus" form:"objStatus"`
	Path        string         `json:"path" form:"path"`
	ParentId    string         `json:"parentId" form:"parentId"`
	Show        bool           `json:"show" form:"show"`
	NotShow     bool           `json:"notShow" form:"notShow"`
	WorkspaceId string         `json:"workspaceId" form:"workspaceId"`
}
