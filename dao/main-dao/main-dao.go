package MainDao

import (
	"errors"
	yiuBool "github.com/fidelyiu/yiu-go-tool/bool"
	yiuStr "github.com/fidelyiu/yiu-go-tool/string"
	"go.etcd.io/bbolt"
	"yiu/yiu-reader/bean"
	EditSoftDao "yiu/yiu-reader/dao/edit-soft-dao"
	WorkspaceDao "yiu/yiu-reader/dao/workspace-dao"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

const tableName = FieldUtil.MainTable

// GetCurrentWorkspaceId 获取当前工作空间Id
func GetCurrentWorkspaceId() (string, error) {
	bytes, err := getMainTableByKey(FieldUtil.CurrentWorkspaceIdField, "无当前工作空间")
	if err != nil {
		return "", err
	}
	result := string(bytes)
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
	return setMainTableByKey(FieldUtil.CurrentWorkspaceIdField, id)
}

func GetMainBoxShowText() (bool, error) {
	bytes, err := getMainTableByKey(FieldUtil.MainBoxShowBtnText, "暂未保存MainBox按钮提示Key")
	if err != nil {
		return true, err
	}
	result := yiuStr.IsTrue(string(bytes))
	return result, nil
}

func SetMainBoxShowText(b bool) error {
	return setMainTableByKey(FieldUtil.MainBoxShowBtnText, yiuBool.ToStr(b))
}

func GetMainBoxShowIcon() (bool, error) {
	bytes, err := getMainTableByKey(FieldUtil.MainBoxShowIcon, "暂未保存MainBox图标Key")
	if err != nil {
		return true, err
	}
	result := yiuStr.IsTrue(string(bytes))
	return result, nil
}

func SetMainBoxShowIcon(b bool) error {
	return setMainTableByKey(FieldUtil.MainBoxShowIcon, yiuBool.ToStr(b))
}

func GetMainBoxShowNum() (bool, error) {
	bytes, err := getMainTableByKey(FieldUtil.MainBoxShowNum, "暂未保存MainBox序号Key")
	if err != nil {
		return true, err
	}
	result := yiuStr.IsTrue(string(bytes))
	return result, nil
}

func SetMainBoxShowNum(b bool) error {
	return setMainTableByKey(FieldUtil.MainBoxShowNum, yiuBool.ToStr(b))
}

func GetSidebarStatus() (bool, error) {
	bytes, err := getMainTableByKey(FieldUtil.SidebarStatusField, "暂未侧边栏状态Key")
	if err != nil {
		return true, err
	}
	result := yiuStr.IsTrue(string(bytes))
	return result, nil
}

func SetSidebarStatus(b bool) error {
	return setMainTableByKey(FieldUtil.SidebarStatusField, yiuBool.ToStr(b))
}

func setMainTableByKey(key, data string) error {
	err := bean.GetDbBean().Update(func(tx *bbolt.Tx) error {
		table := tx.Bucket([]byte(tableName))
		err := table.Put([]byte(key), []byte(data))
		return err
	})
	return err
}

func getMainTableByKey(key, errStr string) ([]byte, error) {
	var result []byte
	err := bean.GetDbBean().View(func(tx *bbolt.Tx) error {
		mainTable := tx.Bucket([]byte(tableName))
		result = mainTable.Get([]byte(key))
		return nil
	})
	if len(result) == 0 {
		return result, errors.New(errStr)
	}
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetEditSoftId 获取当前编辑软件的Id
func GetEditSoftId() (string, error) {
	bytes, err := getMainTableByKey(FieldUtil.CurrentEditId, "无当前编辑软件")
	if err != nil {
		return "", err
	}
	result := string(bytes)
	return result, nil
}

// SetEditSoftId 设置当前编辑软件
func SetEditSoftId(id string) error {
	if id == "" {
		return errors.New("修改当前编辑软件ID失败，修改key不能为空")
	}
	editSoft, err := EditSoftDao.FindById(id)
	if err != nil || editSoft.CheckPath() != nil {
		return errors.New("设置的编辑软件无效")
	}
	return setMainTableByKey(FieldUtil.CurrentEditId, id)
}

func GetNoteTextDocument() (bool, error) {
	bytes, err := getMainTableByKey(FieldUtil.NoteTextDocument, "暂未笔记页面文档文字提示Key")
	if err != nil {
		return true, err
	}
	result := yiuStr.IsTrue(string(bytes))
	return result, nil
}

func SetNoteTextDocument(b bool) error {
	return setMainTableByKey(FieldUtil.NoteTextDocument, yiuBool.ToStr(b))
}

func GetNoteTextMainPoint() (bool, error) {
	bytes, err := getMainTableByKey(FieldUtil.NoteTextMainPoint, "暂未笔记页面大纲文字提示Key")
	if err != nil {
		return true, err
	}
	result := yiuStr.IsTrue(string(bytes))
	return result, nil
}

func SetNoteTextMainPoint(b bool) error {
	return setMainTableByKey(FieldUtil.NoteTextMainPoint, yiuBool.ToStr(b))
}

func GetNoteTextDir() (bool, error) {
	bytes, err := getMainTableByKey(FieldUtil.NoteTextDir, "暂未笔记页面目录文字提示Key")
	if err != nil {
		return true, err
	}
	result := yiuStr.IsTrue(string(bytes))
	return result, nil
}

func SetNoteTextDir(b bool) error {
	return setMainTableByKey(FieldUtil.NoteTextDir, yiuBool.ToStr(b))
}
