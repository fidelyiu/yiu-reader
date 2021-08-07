package entity

import (
	"errors"
	yiuSErr "github.com/fidelyiu/yiu-go-tool/error_s"
	yiuFile "github.com/fidelyiu/yiu-go-tool/file"
	"yiu/yiu-reader/model/enum"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

type Image struct {
	Id     string         `json:"id"`
	Path   string         `json:"path"`   // 图片路径，相对于`/.yiu/image`的路径
	Status enum.ObjStatus `json:"status"` // 状态
	Src    string         `json:"src"`    // 源图片路径
}

func (i *Image) CheckPath() error {
	if !yiuFile.IsExists(FieldUtil.ImageAdd + i.Path) {
		i.Status = enum.ObjStatusInvalid
		return errors.New("图片 '" + i.Path + "' 不是有效文件")
	}
	i.Status = enum.ObjStatusValid
	return nil
}

func (i *Image) Check() error {
	return yiuSErr.ToErrorBySep(" & ",
		i.CheckPath(),
	)
}
