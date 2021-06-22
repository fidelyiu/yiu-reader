package LayoutDao

import (
	"encoding/json"
	"github.com/go-basic/uuid"
	"sort"
	"time"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/entity"
	FieldUtil "yiu/yiu-reader/util/field-util"
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
