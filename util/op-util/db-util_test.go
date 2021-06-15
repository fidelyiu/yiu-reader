package OpUtil

import (
	"path"
	"testing"
)

func TestPathDir(t *testing.T) {
	t.Log(path.Dir(""))
	t.Log(path.Dir("dir"))
	t.Log(path.Dir("/dir"))
	t.Log(path.Dir("./dir"))
	t.Log(path.Dir("./dir/file"))
	t.Log(path.Dir("./dir/subDir"))
}
