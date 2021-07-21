package dto

import "yiu/yiu-reader/model/enum"

type EditSoftSearchDto struct {
	ObjStatus enum.ObjStatus `json:"objStatus" form:"objStatus"`
	Key       string         `json:"key" form:"key"`
	Name      string         `json:"name" form:"name"`
	Path      string         `json:"path" form:"path"`
	SortType  enum.SortType  `json:"sortType" form:"sortType"`
}
