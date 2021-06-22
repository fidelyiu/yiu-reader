package enum

type LayoutType int32

const (
	_                   LayoutType = iota
	LayoutTypeLink                 // 链接
	LayoutTypeDirectory            // 目录
	LayoutTypeFile                 // Markdown文件
	LayoutTypeWorkspace            // 工作空间
	LayoutTypeMainBox              // 主盒子，展示当前工作空间的内容
)
