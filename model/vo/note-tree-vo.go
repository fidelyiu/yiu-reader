package vo

import "yiu/yiu-reader/model/entity"

type NoteTreeVo struct {
	Data  entity.Note
	Child []NoteTreeVo
}
