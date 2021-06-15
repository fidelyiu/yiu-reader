package MainDao

import (
	"errors"
	"go.etcd.io/bbolt"
	"yiu/yiu-reader/bean"
	WorkspaceDao "yiu/yiu-reader/dao/workspace-dao"
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

// SetCurrentWorkspaceId 设置当前工作空间
func SetCurrentWorkspaceId(id string) error {
	if id == "" {
		return errors.New("修改当前工作空间ID失败，修改key不能为空")
	}
	workspace, err := WorkspaceDao.FindById(id)
	if err != nil || workspace.CheckPath() != nil {
		return errors.New("设置的工作空间无效")
	}
	err = bean.GetDbBean().Update(func(tx *bbolt.Tx) error {
		table := tx.Bucket([]byte(tableName))
		err := table.Put([]byte(FieldUtil.CurrentWorkspaceIdField), []byte(id))
		return err
	})
	return err
}
