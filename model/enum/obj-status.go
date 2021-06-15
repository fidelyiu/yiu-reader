package enum

type ObjStatus int32

const (
	ObjStatusNoValue ObjStatus = iota // 无状态
	ObjStatusInvalid                  // 无效
	ObjStatusValid                    // 有效
)
