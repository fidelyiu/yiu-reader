package vo

type DbSearchVo struct {
	Data        []DbSearchItemVo `json:"data"`
	TotalRow    int              `json:"totalRow"`
	TotalPage   int              `json:"totalPage"`
	CurrentPage int              `json:"currentPage"`
	PageSize    int              `json:"pageSize"`
	HasPrevious bool             `json:"hasPrevious"`
	HasNext     bool             `json:"hasNext"`
}
