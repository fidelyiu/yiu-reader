package WorkspaceDao

import (
	"encoding/json"
	"errors"
	YiuStr "github.com/fidelyiu/yiu-go/string"
	"github.com/go-basic/uuid"
	"sort"
	"strings"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/dto"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

const (
	tableName  = FieldUtil.WorkspaceTable
	entityName = "工作空间"
)

func FindById(id string) (entity.Workspace, error) {
	v, err := dao.FindByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
	if len(v) == 0 {
		return entity.Workspace{}, errors.New("[" + id + "]" + entityName + "无效")
	}
	if err != nil {
		return entity.Workspace{}, err
	}
	var result entity.Workspace
	err = json.Unmarshal([]byte(v), &result)
	return result, err
}

func Save(entity *entity.Workspace) error {
	entity.Id = uuid.New()
	totalNum, err := dao.CountAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return err
	}
	entity.SortNum = totalNum + 1
	buf, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return dao.SaveByTableNameAndKey(bean.GetDbBean(), tableName, entity.Id, buf, entityName)
}

func FindAllBySearchDto(dto dto.WorkspaceSearchDto) ([]entity.Workspace, error) {
	stringList, err := dao.FindAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Workspace, 0)
	for _, v := range stringList {
		var resultItem entity.Workspace
		err := json.Unmarshal([]byte(v), &resultItem)
		if err == nil {
			appendItem := true
			// 关键字过滤
			if YiuStr.IsNotBlank(dto.Key) &&
				!strings.Contains(resultItem.Name, dto.Key) &&
				!strings.Contains(resultItem.Name, dto.Key) {
				appendItem = false
			}
			statusErr := resultItem.CheckPath()
			// 名称过滤
			if YiuStr.IsNotBlank(dto.Name) && !strings.Contains(resultItem.Name, dto.Name) {
				appendItem = false
			}
			// 路径过滤
			if YiuStr.IsNotBlank(dto.Path) && !strings.Contains(resultItem.Path, dto.Path) {
				appendItem = false
			}
			// 过滤是否有效
			switch dto.ObjStatus {
			case enum.ObjStatusValid:
				if statusErr != nil {
					appendItem = false
				}
			case enum.ObjStatusInvalid:
				if statusErr == nil {
					appendItem = false
				}
			}
			if appendItem {
				result = append(result, resultItem)
			}
		}
	}
	return result, err
}

func SearchByDto(dto dto.WorkspaceSearchDto) ([]entity.Workspace, error) {
	entityList, err := FindAllBySearchDto(dto)
	if err != nil {
		return nil, err
	}
	if dto.SortType == enum.ASE {
		sort.Slice(entityList, func(i, j int) bool {
			return entityList[i].SortNum > entityList[j].SortNum
		})
	} else {
		sort.Slice(entityList, func(i, j int) bool {
			return entityList[i].SortNum < entityList[j].SortNum
		})
	}
	return entityList, err
}

func Update(updateEntity *entity.Workspace) error {
	var dbEntity entity.Workspace
	if YiuStr.IsNotBlank(updateEntity.Id) {
		dbEntity, _ = FindById(updateEntity.Id)
	}
	// 不能修改序号
	if YiuStr.IsBlank(dbEntity.Id) {
		return errors.New("修改" + entityName + "报错，id不能为空")
	} else {
		updateEntity.SortNum = dbEntity.SortNum
	}
	buf, err := json.Marshal(updateEntity)
	if err != nil {
		return err
	}
	return dao.UpdateByTableNameAndKey(bean.GetDbBean(), tableName, updateEntity.Id, buf, entityName)
}

func DeleteById(id string) error {
	return dao.DeleteByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
}
