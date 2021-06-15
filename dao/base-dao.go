package dao

import (
	"errors"
	YiuStr "github.com/fidelyiu/yiu-go/string"
	YiuStrList "github.com/fidelyiu/yiu-go/string_list"
	"go.etcd.io/bbolt"
)

// FindAllByTableName 根据表明查找所有数据
func FindAllByTableName(db *bbolt.DB, tableName string) ([]string, error) {
	result := make([]string, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		table := GetTableByName(tx, tableName)
		err := table.ForEach(func(k, v []byte) error {
			result = append(result, string(v))
			return nil
		})
		return err
	})
	return result, err
}

// CountAllByTableName 根据表明统计所有数据
func CountAllByTableName(db *bbolt.DB, tableName string) (int, error) {
	var result int
	err := db.View(func(tx *bbolt.Tx) error {
		table := GetTableByName(tx, tableName)
		err := table.ForEach(func(_, _ []byte) error {
			result++
			return nil
		})
		return err
	})
	return result, err
}

// FindByTableNameAndKey 根据表名&Key查找一个数据
func FindByTableNameAndKey(db *bbolt.DB, tableName string, key string, entityName string) (string, error) {
	if key == "" {
		return "", errors.New("查询" + entityName + "报错，key不能为空")
	}
	var result string
	err := db.View(func(tx *bbolt.Tx) error {
		table := GetTableByName(tx, tableName)
		result = string(table.Get([]byte(key)))
		return nil
	})
	return result, err
}

// IsEffectiveByTableNameAndKey 判断该ID是否有效
func IsEffectiveByTableNameAndKey(db *bbolt.DB, tableName string, key string, entityName string) bool {
	v, err := FindByTableNameAndKey(db, tableName, key, entityName)
	if err != nil || len(v) == 0 {
		return false
	}
	return true
}

// SaveByTableNameAndKey 根据表明保存一条数据，key&数据不能为空
func SaveByTableNameAndKey(db *bbolt.DB, tableName string, key string, entityByte []byte, entityName string) error {
	if key == "" {
		return errors.New("保存" + entityName + "报错，key不能为空")
	}
	if len(entityByte) == 0 {
		return errors.New("保存" + entityName + "报错，数据不能为空")
	}
	err := db.Update(func(tx *bbolt.Tx) error {
		table := GetTableByName(tx, tableName)
		err := table.Put([]byte(key), entityByte)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateByTableNameAndKey 根据表名、Key、名称修改一个数据，如果数据不存在报错
func UpdateByTableNameAndKey(db *bbolt.DB, tableName string, key string, entityByte []byte, entityName string) error {
	dbEntity, err := FindByTableNameAndKey(db, tableName, key, entityName)
	if err != nil {
		return err
	}
	if dbEntity == "" {
		return errors.New("修改" + entityName + "报错，数据不存在")
	}
	if key == "" {
		return errors.New("修改" + entityName + "报错，key不能为空")
	}
	if len(entityByte) == 0 {
		return errors.New("修改" + entityName + "报错，数据不能为空")
	}
	return SaveByTableNameAndKey(db, tableName, key, entityByte, entityName)
}

// DeleteByTableNameAndKey 根据表名、Key删除一条数据
func DeleteByTableNameAndKey(db *bbolt.DB, tableName string, key string, entityName string) error {
	if key == "" {
		return errors.New("删除" + entityName + "报错，key不能为空")
	}
	err := db.Update(func(tx *bbolt.Tx) error {
		table := GetTableByName(tx, tableName)
		err := table.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetTableByName 根据字符串获取数据库中的表
func GetTableByName(tx *bbolt.Tx, tableName string) *bbolt.Bucket {
	if YiuStr.IsBlank(tableName) {
		return nil
	}
	tableNameList := YiuStr.ToStrList(tableName, ".")
	baseBucket := tx.Bucket([]byte(tableNameList[0]))
	if len(tableNameList) <= 1 {
		return baseBucket
	}
	YiuStrList.OpDeleteByIndex(&tableNameList, 0)
	for _, v := range tableNameList {
		baseBucket = baseBucket.Bucket([]byte(v))
	}
	return baseBucket
}

func GetCustomizeWork(db *bbolt.DB, opFunc func(tx *bbolt.Tx) error) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer func(tx *bbolt.Tx) {
		_ = tx.Rollback()
	}(tx)
	err = opFunc(tx)
	if err != nil {
		return err
	}
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
