package enum

type NoteReadResult int32

const (
	NoteReadResultFail = iota
	NoteReadResultImport
	NoteReadResultNotImport
)
