package vo

import "yiu/yiu-reader/model/enum"

type WorkspaceRefreshVo struct {
	Message string                    `json:"message"`
	Type    enum.WorkspaceRefreshType `json:"type"`
	Result  interface{}               `json:"result"`
}
