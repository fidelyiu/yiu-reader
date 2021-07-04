package vo

import "yiu/yiu-reader/model/enum"

// NoteReadVo 读取本文件时，结束的VO
type NoteReadVo struct {
	Path   string              // 路径
	Result enum.NoteReadResult // 读取结果
}
