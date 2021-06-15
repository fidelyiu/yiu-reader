package enum

type ChangeSortType int32

const (
	ChangeSortTypeUp   ChangeSortType = iota // 升序
	ChangeSortTypeDown                       // 降序
)
