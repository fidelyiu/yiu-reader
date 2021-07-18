package PathUtil

import (
	yiuDir "github.com/fidelyiu/yiu-go-tool/dir"
	yiuFile "github.com/fidelyiu/yiu-go-tool/file"
	"path"
	"path/filepath"
)

// IsValidDir 是否是有效的实体目录路径
// 1. 需要是绝对路径
// 2. 路径本身有效
// 3. 文件夹
func IsValidDir(str string) bool {
	return filepath.IsAbs(str) && yiuDir.IsExists(str)
}

// IsValidFile 是否是有效的实体文件路径
// 1. 需要是绝对路径
// 2. 路径本身有效
// 3. 文件
func IsValidFile(str string) bool {
	return filepath.IsAbs(str) && yiuFile.IsExists(str)
}

// IsValidMarkdown 是否是有效的实体Markdown文件路径
// 1. 需要是绝对路径
// 2. 路径本身有效
// 3. 文件以`.md`结尾
func IsValidMarkdown(str string) bool {
	return path.IsAbs(str) && yiuFile.IsExists(str) && path.Ext(str) == ".md"
}
