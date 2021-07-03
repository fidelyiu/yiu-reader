package enum

type DefinitionStatus int32

const (
	DefinitionStatusNoValue       DefinitionStatus = iota // 未定义
	DefinitionStatusHasDefinition                         // 已定义
)
