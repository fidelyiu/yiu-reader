package LayoutDao

import (
	"encoding/json"
	"errors"
	"github.com/go-basic/uuid"
	"go.etcd.io/bbolt"
	"sort"
	"time"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/entity"
	FieldUtil "yiu/yiu-reader/util/field-util"
	LayoutUtil "yiu/yiu-reader/util/layout-util"
)

const (
	tableName  = FieldUtil.LayoutTable
	entityName = "布局"
)

func Save(entity *entity.Layout) error {
	entity.Id = uuid.New()
	entity.UpdateTime = time.Now()
	buf, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return dao.SaveByTableNameAndKey(bean.GetDbBean(), tableName, entity.Id, buf, entityName)
}

func FindAll() ([]entity.Layout, error) {
	stringList, err := dao.FindAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Layout, 0)
	for _, v := range stringList {
		var resultItem entity.Layout
		err := json.Unmarshal([]byte(v), &resultItem)
		if err == nil {
			result = append(result, resultItem)
		}
	}
	return result, err
}

func FindAllSortByUpdateTime() ([]entity.Layout, error) {
	entityList, err := FindAll()
	if err != nil {
		return nil, err
	}
	sort.Slice(entityList, func(i, j int) bool {
		return entityList[i].UpdateTime.After(entityList[j].UpdateTime)
	})
	return entityList, err
}

func DeleteById(id string) error {
	return dao.DeleteByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
}

func FindById(id string) (entity.Layout, error) {
	v, err := dao.FindByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
	if len(v) == 0 {
		return entity.Layout{}, errors.New("[" + id + "]" + entityName + "无效")
	}
	if err != nil {
		return entity.Layout{}, err
	}
	var result entity.Layout
	err = json.Unmarshal([]byte(v), &result)
	return result, err
}

func Update(entity *entity.Layout) error {
	entity.UpdateTime = time.Now()
	buf, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return dao.UpdateByTableNameAndKey(bean.GetDbBean(), tableName, entity.Id, buf, entityName)
}

func FormatAll(maxX int) error {
	// 开始事务
	tx, err := bean.GetDbBean().Begin(true)
	if err != nil {
		return err
	}
	// 结尾回滚事务
	defer func(tx *bbolt.Tx) {
		_ = tx.Rollback()
	}(tx)
	// 获取表
	table := dao.GetTableByName(tx, tableName)
	// 获取所有数据
	layoutList := make([]entity.Layout, 0)
	err = table.ForEach(func(k, v []byte) error {
		var vItem entity.Layout
		err := json.Unmarshal(v, &vItem)
		if err != nil {
			return err
		}
		layoutList = append(layoutList, vItem)
		return nil
	})
	if err != nil {
		return err
	}

	// 时间排序
	sort.Slice(layoutList, func(i, j int) bool {
		return layoutList[i].UpdateTime.After(layoutList[j].UpdateTime)
	})

	// 格式化
	invalidLayoutList := LayoutUtil.OutInvalidLayout(layoutList, maxX)

	// 无效布局全部更新一遍
	for _, v := range invalidLayoutList {
		buf, err := json.Marshal(v)
		if err != nil {
			continue
		}
		err = table.Put([]byte(v.Id), buf)
	}

	// 提交事务
	err = tx.Commit()
	return err
}
