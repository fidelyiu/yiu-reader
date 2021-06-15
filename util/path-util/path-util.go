package PathUtil

import (
	YiuDir "github.com/fidelyiu/yiu-go/dir"
	YiuFile "github.com/fidelyiu/yiu-go/file"
	"path"
	"path/filepath"
)

// IsValidDir 是否是有效的实体目录路径
// 1. 需要是绝对路径
// 2. 路径本身有效
func IsValidDir(str string) bool {
	return filepath.IsAbs(str) && YiuDir.IsExists(str)
}

// IsValidMarkdown 是否是有效的实体Markdown文件路径
// 1. 需要是绝对路径
// 2. 路径本身有效
// 3. 文件以`.md`结尾
func IsValidMarkdown(str string) bool {
	return path.IsAbs(str) && YiuFile.IsExists(str) && path.Ext(str) == ".md"
}
