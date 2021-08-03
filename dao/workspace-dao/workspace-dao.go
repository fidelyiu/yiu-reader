package WorkspaceDao

import (
	"encoding/json"
	"errors"
	yiuStr "github.com/fidelyiu/yiu-go-tool/string"
	"github.com/go-basic/uuid"
	"go.etcd.io/bbolt"
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

func Save(saveEntity *entity.Workspace) error {
	saveEntity.Id = uuid.New()
	entityStrList, err := dao.FindAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return err
	}

	pathHasInsert := false
	containItem := false
	containEdItem := false
	errName := ""
	// 查看路径是否已经保存过
	for _, v := range entityStrList {
		var entityItem entity.Workspace
		err = json.Unmarshal([]byte(v), &entityItem)
		if err != nil {
			continue
		}
		_ = entityItem.CheckPath()
		if entityItem.Status != enum.ObjStatusValid {
			continue
		}

		if entityItem.Path == saveEntity.Path {
			pathHasInsert = true
			errName = entityItem.Name
			break
		}
		if strings.Contains(saveEntity.Path, entityItem.Path) {
			containItem = true
			errName = entityItem.Name
			break
		}
		if strings.Contains(entityItem.Path, saveEntity.Path) {
			containEdItem = true
			errName = entityItem.Name
			break
		}
	}
	if pathHasInsert {
		return errors.New("该路径已保存为'" + errName + "'，请勿重复保存工作空间")
	}
	if containItem {
		return errors.New("该路径包含'" + errName + "'工作空间，请勿交叉保存工作空间，这会使文档所属无法定位。")
	}
	if containEdItem {
		return errors.New("该路径被包含在'" + errName + "'工作空间中，请勿交叉保存工作空间，这会使文档所属无法定位。")
	}

	saveEntity.SortNum = len(entityStrList) + 1
	buf, err := json.Marshal(saveEntity)
	if err != nil {
		return err
	}
	return dao.SaveByTableNameAndKey(bean.GetDbBean(), tableName, saveEntity.Id, buf, entityName)
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
			if yiuStr.IsNotBlank(dto.Key) &&
				!strings.Contains(resultItem.Name, dto.Key) &&
				!strings.Contains(resultItem.Path, dto.Key) {
				appendItem = false
			}
			statusErr := resultItem.CheckPath()
			// 名称过滤
			if yiuStr.IsNotBlank(dto.Name) && !strings.Contains(resultItem.Name, dto.Name) {
				appendItem = false
			}
			// 路径过滤
			if yiuStr.IsNotBlank(dto.Path) && !strings.Contains(resultItem.Path, dto.Path) {
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
	if dto.SortType == enum.SortTypeAse {
		sort.Slice(entityList, func(i, j int) bool {
			return entityList[i].SortNum < entityList[j].SortNum
		})
	} else {
		sort.Slice(entityList, func(i, j int) bool {
			return entityList[i].SortNum > entityList[j].SortNum
		})
	}
	return entityList, err
}

func Update(updateEntity *entity.Workspace) error {
	var dbEntity entity.Workspace
	if yiuStr.IsNotBlank(updateEntity.Id) {
		dbEntity, _ = FindById(updateEntity.Id)
	}
	// 不能修改序号
	if yiuStr.IsBlank(dbEntity.Id) {
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
	// return dao.DeleteByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
	// 开事务删笔记
	// 检查ID是否有效
	if !dao.IsEffectiveByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName) {
		return errors.New("修改" + entityName + "报错，id无效")
	}
	// 开始事务
	tx, err := bean.GetDbBean().Begin(true)
	if err != nil {
		return err
	}
	// 结尾回滚事务
	defer func(tx *bbolt.Tx) {
		_ = tx.Rollback()
	}(tx)
	noteTable := dao.GetTableByName(tx, FieldUtil.NoteTable)

	var noteList []entity.Note
	err = noteTable.ForEach(func(k, v []byte) error {
		var vItem entity.Note
		itemErr := json.Unmarshal(v, &vItem)
		if itemErr != nil {
			return itemErr
		}
		if vItem.WorkspaceId == id {
			noteList = append(noteList, vItem)
		}
		return nil
	})

	for i := range noteList {
		delErr := noteTable.Delete([]byte(noteList[i].Id))
		if delErr != nil {
			return delErr
		}
	}

	// 获取表
	table := dao.GetTableByName(tx, tableName)

	err = table.Delete([]byte(id))
	if err != nil {
		return err
	}

	// 提交事务
	err = tx.Commit()
	return err
}

func ChangeSort(id string, changeType enum.ChangeSortType) error {
	// 检查ID是否有效
	if !dao.IsEffectiveByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName) {
		return errors.New("修改" + entityName + "报错，id无效")
	}
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
	workspaceList := make([]entity.Workspace, 0)
	err = table.ForEach(func(k, v []byte) error {
		var vItem entity.Workspace
		itemErr := json.Unmarshal(v, &vItem)
		if itemErr != nil {
			return itemErr
		}
		workspaceList = append(workspaceList, vItem)
		return nil
	})
	if err != nil {
		return err
	}

	// 排序
	sort.Slice(workspaceList, func(i, j int) bool {
		return workspaceList[i].SortNum < workspaceList[j].SortNum
	})

	targetIndex := -1
	for i, v := range workspaceList {
		if v.Id == id {
			targetIndex = i
		}
		workspaceList[i].SortNum = i + 1
	}
	if targetIndex == -1 {
		return errors.New("查找" + entityName + "报错，id无效")
	}

	if changeType == enum.ChangeSortTypeUp {
		if targetIndex-1 >= 0 {
			workspaceList[targetIndex].SortNum--
			workspaceList[targetIndex-1].SortNum++
		}
	} else {
		if targetIndex+1 < len(workspaceList) {
			workspaceList[targetIndex].SortNum++
			workspaceList[targetIndex+1].SortNum--
		}
	}

	// 遍历修改序号
	for _, v := range workspaceList {
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
