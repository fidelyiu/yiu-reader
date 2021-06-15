package enum

type ResultType string

const (
	ResultTypeSuccess ResultType = "success" // 成功
	ResultTypeError   ResultType = "error"   // 失败
	ResultTypeWarning ResultType = "warning" // 警告
)
