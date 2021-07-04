package enum

type NoteReadResult int32

const (
	NoteReadResultFail      = iota // 失败
	NoteReadResultStart            // 开始
	NoteReadResultImport           // 已导入
	NoteReadResultNotImport        // 未导入
)
