package entity

import (
	"errors"
	YiuErrorList "github.com/fidelyiu/yiu-go/error_list"
	"time"
	"yiu/yiu-reader/model/enum"
)

type Layout struct {
	Id         string                 `json:"id"`         // Uuid
	Type       enum.LayoutType        `json:"type"`       // 类型
	Status     enum.ObjStatus         `json:"status"`     // 状态
	Width      int                    `json:"width"`      // 宽度
	Height     int                    `json:"height"`     // 高度
	Left       int                    `json:"left"`       // 距离左边
	Top        int                    `json:"top"`        // 距离顶部
	Setting    map[string]interface{} `json:"setting"`    // 设置，根据类型而定
	UpdateTime time.Time              `json:"updateTime"` // 最后更新时间
}

func (l *Layout) CheckType() error {
	if l.Type <= enum.LayoutTypeLink-1 || l.Type >= enum.LayoutTypeMainBox+1 {
		return errors.New("布局[Type]无效")
	}
	return nil
}

func (l *Layout) CheckStatus() error {
	if l.Status <= enum.ObjStatusInvalid-1 || l.Status >= enum.ObjStatusValid+1 {
		return errors.New("布局状态无效")
	}
	return nil
}

func (l *Layout) CheckWidth() error {
	if l.Width <= 0 {
		return errors.New("布局宽度无效")
	}
	return nil
}

func (l *Layout) CheckHeight() error {
	if l.Width <= 0 {
		return errors.New("布局高度无效")
	}
	return nil
}

func (l *Layout) CheckSetting() error {
	switch l.Type {
	case enum.LayoutTypeLink:
		if l.Setting["name"] == nil {
			return errors.New("链接名称不能为空")
		}
		if l.Setting["url"] == nil {
			return errors.New("链接地址不能为空")
		}
	}
	return nil
}

func (l *Layout) Check() error {
	return YiuErrorList.ToError([]error{
		l.CheckType(),
		l.CheckStatus(),
		l.CheckWidth(),
		l.CheckHeight(),
	}, " & ")
}
