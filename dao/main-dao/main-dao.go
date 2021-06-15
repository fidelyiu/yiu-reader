package MainDao

import (
	"errors"
	"go.etcd.io/bbolt"
	"yiu/yiu-reader/bean"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

const tableName = FieldUtil.MainTable

// GetCurrentWorkspaceId 获取当前工作空间Id
func GetCurrentWorkspaceId() (string, error) {
	result := ""
	err := bean.GetDbBean().View(func(tx *bbolt.Tx) error {
		mainTable := tx.Bucket([]byte(tableName))
		v := mainTable.Get([]byte(FieldUtil.CurrentWorkspaceIdField))
		result = string(v)
		return nil
	})
	if result == "" {
		return "", errors.New("无当前工作空间")
	}
	if err != nil {
		return "", err
	}
	return result, nil
}
