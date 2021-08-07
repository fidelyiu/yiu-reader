package entity

import (
	"errors"
	yiuSErr "github.com/fidelyiu/yiu-go-tool/error_s"
	yiuFile "github.com/fidelyiu/yiu-go-tool/file"
	"yiu/yiu-reader/model/enum"
	FieldUtil "yiu/yiu-reader/util/field-util"
	PathUtil "yiu/yiu-reader/util/path-util"
)

type Image struct {
	Id        string         `json:"id"`
	Path      string         `json:"path"`      // 图片路径，相对于`/.yiu/image`的路径
	Status    enum.ObjStatus `json:"status"`    // 状态
	Src       string         `json:"src"`       // 源图片路径
	SrcStatus enum.ObjStatus `json:"srcStatus"` // 源状态
	IsUpload  bool           `json:"isUpload"`  // 是否是自己上传的
}

func (i *Image) CheckPath() error {
	if !yiuFile.IsExists(FieldUtil.ImageAdd + i.Path) {
		i.Status = enum.ObjStatusInvalid
		return errors.New("图片 '" + i.Path + "' 不是有效文件")
	}
	i.Status = enum.ObjStatusValid
	return nil
}

func (i *Image) CheckSrc() error {
	if !PathUtil.IsValidFile(i.Src) {
		i.SrcStatus = enum.ObjStatusInvalid
		return errors.New("图片 '" + i.Src + "' 不是有效文件的绝对路径")
	}
	i.SrcStatus = enum.ObjStatusValid
	return nil
}

func (i *Image) Check() error {
	return yiuSErr.ToErrorBySep(" & ",
		i.CheckPath(),
		i.CheckSrc(),
	)
}
