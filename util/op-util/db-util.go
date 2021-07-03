package OpUtil

import (
	YiuDir "github.com/fidelyiu/yiu-go/dir"
	YiuLogger "github.com/fidelyiu/yiu-go/logger"
	"go.etcd.io/bbolt"
	"path"
	"yiu/yiu-reader/bean"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

func CreateDB(path string) {
	db, err := OpenBoltDB(".yiu/yiu-reader.db")
	if err != nil {
		YiuLogger.LogErrorLn("打开数据库出错：")
		YiuLogger.LogErrorLn(err)
		return
	}
	bean.SetDbBean(db)
}

func CloseDB() {
	db := bean.GetDbBean()
	if db == nil {
		return
	}
	err := db.Close()
	if err != nil {
		YiuLogger.LogErrorLn("关闭数据库出错：")
		YiuLogger.LogErrorLn(err)
		return
	}
}

func OpenBoltDB(dbPath string) (*bbolt.DB, error) {
	dirPath := path.Dir(dbPath)
	if !YiuDir.IsExists(dirPath) {
		err := YiuDir.OpMkDirAll(dirPath)
		if err != nil {
			return nil, err
		}
	}
	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		// 检查创建MainTable
		_, createErr := tx.CreateBucketIfNotExists([]byte(FieldUtil.MainTable))
		if createErr != nil {
			return createErr
		}
		_, createErr = tx.CreateBucketIfNotExists([]byte(FieldUtil.LayoutTable))
		if createErr != nil {
			return createErr
		}
		_, createErr = tx.CreateBucketIfNotExists([]byte(FieldUtil.WorkspaceTable))
		if createErr != nil {
			return createErr
		}
		_, createErr = tx.CreateBucketIfNotExists([]byte(FieldUtil.NoteTable))
		if createErr != nil {
			return createErr
		}
		_, createErr = tx.CreateBucketIfNotExists([]byte(FieldUtil.MarkdownTable))
		if createErr != nil {
			return createErr
		}
		_, createErr = tx.CreateBucketIfNotExists([]byte(FieldUtil.ImageCacheTable))
		if createErr != nil {
			return createErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
