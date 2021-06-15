package response

import "yiu/yiu-reader/model/enum"

type YiuReaderResponse struct {
	Code    int             `json:"code"`
	Type    enum.ResultType `json:"type"`
	Result  interface{}     `json:"result"`
	Message string          `json:"message"`
}

func (i *YiuReaderResponse) SetCode(code int) *YiuReaderResponse {
	i.Code = code
	return i
}
func (i *YiuReaderResponse) SetType(resultType enum.ResultType) *YiuReaderResponse {
	i.Type = resultType
	return i
}
func (i *YiuReaderResponse) SetResult(result interface{}) *YiuReaderResponse {
	i.Result = result
	return i
}
func (i *YiuReaderResponse) SetMessage(message string) *YiuReaderResponse {
	i.Message = message
	return i
}
func (i *YiuReaderResponse) ToError(message string) {
	i.Message = message
	i.Type = enum.ResultTypeError
}
