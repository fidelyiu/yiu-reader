package NoteDao

import (
	"encoding/json"
	"errors"
	"go.etcd.io/bbolt"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/entity"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

const (
	tableName  = FieldUtil.NoteTable
	entityName = "笔记"
)

func FindAll() ([]entity.Note, error) {
	stringList, err := dao.FindAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Note, 0)
	for _, v := range stringList {
		var resultItem entity.Note
		err := json.Unmarshal([]byte(v), &resultItem)
		if err == nil {
			result = append(result, resultItem)
		}
	}
	return result, err
}

// FindByAbsPath 根据绝对路径找笔记
func FindByAbsPath(absPath string) (entity.Note, error) {
	result := entity.Note{}
	hasData := false
	err := bean.GetDbBean().View(func(tx *bbolt.Tx) error {
		table := dao.GetTableByName(tx, tableName)
		err := table.ForEach(func(_, i []byte) error {
			var tempItem entity.Note
			err := json.Unmarshal(i, &tempItem)
			if err != nil {
				return err
			}
			if tempItem.AbsPath == absPath {
				result = tempItem
				hasData = true
			}
			return nil
		})
		return err
	})
	if hasData {
		return result, err
	} else {
		return result, errors.New("未保存该绝对路径笔记")
	}
}

func SaveAll(noteList []entity.Note) error {
	if len(noteList) == 0 {
		return nil
	}
	err := bean.GetDbBean().Batch(func(tx *bbolt.Tx) error {
		table := dao.GetTableByName(tx, tableName)
		for i := range noteList {
			buf, err := json.Marshal(noteList[i])
			if err != nil {
				return err
			}
			if noteList[i].Id == "" {
				return errors.New("保存" + entityName + "报错，id 不能为空")
			}
			err = table.Put([]byte(noteList[i].Id), buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func Update(entity *entity.Note) error {
	buf, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return dao.UpdateByTableNameAndKey(bean.GetDbBean(), tableName, entity.Id, buf, entityName)
}
