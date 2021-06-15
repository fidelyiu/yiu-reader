package dto

import "yiu/yiu-reader/model/enum"

type WorkspaceSearchDto struct {
	ObjStatus enum.ObjStatus `json:"objStatus"`
	Key       string         `json:"key"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	SortType  enum.SortType  `json:"sortType"`
}
