package NoteDao

import (
	"encoding/json"
	"errors"
	yiuStr "github.com/fidelyiu/yiu-go-tool/string"
	"github.com/go-basic/uuid"
	"go.etcd.io/bbolt"
	"os"
	"sort"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/dto"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	"yiu/yiu-reader/model/vo"
	FieldUtil "yiu/yiu-reader/util/field-util"
	NoteUtil "yiu/yiu-reader/util/note-util"
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
			_ = resultItem.CheckPath()
			result = append(result, resultItem)
		}
	}
	return result, err
}

func FindByWorkspaceId(id string) ([]entity.Note, error) {
	if id == "" {
		return nil, errors.New("id不能为空")
	}
	stringList, err := dao.FindAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Note, 0)
	for _, v := range stringList {
		var resultItem entity.Note
		err := json.Unmarshal([]byte(v), &resultItem)
		if err == nil {
			appendItem := true
			_ = resultItem.CheckPath()
			if id != resultItem.WorkspaceId {
				appendItem = false
			}
			if appendItem {
				result = append(result, resultItem)
			}
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
		_ = result.CheckPath()
		return result, err
	} else {
		return result, errors.New("未保存该绝对路径笔记")
	}
}

func FindByParentId(id string) ([]entity.Note, error) {
	if id == "" {
		return nil, errors.New("id不能为空")
	}
	stringList, err := dao.FindAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Note, 0)
	for _, v := range stringList {
		var resultItem entity.Note
		err := json.Unmarshal([]byte(v), &resultItem)
		if err == nil {
			appendItem := true
			_ = resultItem.CheckPath()
			if id != resultItem.ParentId {
				appendItem = false
			}
			if appendItem {
				result = append(result, resultItem)
			}
		}
	}
	return result, err
}

// FindBySearchDto 根据搜索条件查询
func FindBySearchDto(dto dto.NoteSearchDto) ([]entity.Note, error) {
	stringList, err := dao.FindAllByTableName(bean.GetDbBean(), tableName)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Note, 0)
	for _, v := range stringList {
		var resultItem entity.Note
		err := json.Unmarshal([]byte(v), &resultItem)
		if err == nil {
			appendItem := true
			_ = resultItem.CheckPath()

			// 状态
			if dto.ObjStatus != enum.ObjStatusNoValue && dto.ObjStatus != resultItem.Status {
				appendItem = false
			}

			// 工作空间ID
			if appendItem && yiuStr.IsNotBlank(dto.WorkspaceId) && resultItem.WorkspaceId != dto.WorkspaceId {
				appendItem = false
			}
			// 父路径
			if appendItem && yiuStr.IsNotBlank(dto.ParentId) && resultItem.ParentId != dto.ParentId {
				appendItem = false
			}

			if appendItem && dto.Level != 0 && resultItem.Level != dto.Level {
				appendItem = false
			}

			if appendItem && dto.Show && !resultItem.Show {
				appendItem = false
			}
			if appendItem && dto.NotShow && resultItem.Show {
				appendItem = false
			}

			if appendItem {
				result = append(result, resultItem)
			}
		}
	}
	return result, err
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

func UpdateBatch(list []entity.Note) error {
	return bean.GetDbBean().Batch(func(tx *bbolt.Tx) error {
		table := dao.GetTableByName(tx, tableName)
		for i := range list {
			if list[i].Id == "" {
				continue
			}
			buf, err := json.Marshal(list[i])
			if err != nil {
				return err
			}
			err = table.Put([]byte(list[i].Id), buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func DeleteById(id string) error {
	return dao.DeleteByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
}

func DeleteByIds(ids []string) error {
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
	for i := range ids {
		err = table.Delete([]byte(ids[i]))
	}
	if err != nil {
		return err
	}
	// 提交事务
	err = tx.Commit()
	return err
}

func DeleteDeepById(id string) error {
	return DeleteDeepByIds([]string{id})
}

func DeleteDeepByIds(ids []string) error {
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
	for i := range ids {
		err = deleteDeepById(ids[i], table)
	}
	if err != nil {
		return err
	}
	// 提交事务
	err = tx.Commit()
	return err
}

func deleteDeepById(id string, table *bbolt.Bucket) error {
	var target entity.Note
	err := json.Unmarshal(table.Get([]byte(id)), &target)
	if err != nil {
		return err
	}
	child := make([]entity.Note, 0)
	err = table.ForEach(func(_, v []byte) error {
		var item entity.Note
		err = json.Unmarshal(v, &item)
		if err != nil {
			return err
		}
		if item.ParentId == id {
			child = append(child, item)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(child) != 0 {
		for i := range child {
			err = deleteDeepById(child[i].Id, table)
			if err != nil {
				return err
			}
		}
	}
	err = table.Delete([]byte(id))
	return err
}

func FindById(id string) (entity.Note, error) {
	v, err := dao.FindByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
	if len(v) == 0 {
		return entity.Note{}, errors.New("[" + id + "]" + entityName + "无效")
	}
	if err != nil {
		return entity.Note{}, err
	}
	var result entity.Note
	err = json.Unmarshal([]byte(v), &result)
	return result, err
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

	target, err := FindById(id)
	if err != nil {
		return err
	}

	var brotherList []entity.Note

	err = table.ForEach(func(k, v []byte) error {
		var vItem entity.Note
		itemErr := json.Unmarshal(v, &vItem)
		if itemErr != nil {
			return itemErr
		}
		if vItem.ParentId == target.ParentId &&
			vItem.Level == target.Level &&
			vItem.Show == target.Show {
			brotherList = append(brotherList, vItem)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(brotherList) == 0 {
		return nil
	}

	// 排序
	sort.Slice(brotherList, func(i, j int) bool {
		if brotherList[j].SortNum != 0 {
			return false
		}
		if brotherList[i].SortNum == 0 {
			return brotherList[i].Id == target.Id
		}
		if brotherList[j].SortNum == 0 {
			return true
		}
		return brotherList[i].SortNum < brotherList[j].SortNum
	})

	// 找索引
	targetIndex := -1

	for i := range brotherList {
		if target.Id == brotherList[i].Id {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return errors.New("移动目标异常")
	}

	if target.SortNum == 0 {
		// 0交换，因为上面的排序，目标在0的第一位
		brotherList[targetIndex].SortNum = 130
	} else {
		if len(brotherList) == 1 {
			brotherList[targetIndex].SortNum = 130
		} else {
			// 非0交换
			if changeType == enum.ChangeSortTypeUp {
				if targetIndex != 0 {
					brotherList[targetIndex-1], brotherList[targetIndex] = brotherList[targetIndex], brotherList[targetIndex-1]
				}
			} else {
				if targetIndex != len(brotherList)-1 && brotherList[targetIndex+1].SortNum != 0 {
					brotherList[targetIndex], brotherList[targetIndex+1] = brotherList[targetIndex+1], brotherList[targetIndex]
				}
			}
		}
	}

	for i := range brotherList {
		if brotherList[i].SortNum == 0 {
			break
		}
		brotherList[i].SortNum = i + 1
		updateByte, jsonErr := json.Marshal(brotherList[i])
		if jsonErr != nil {
			return jsonErr
		}
		putErr := table.Put([]byte(brotherList[i].Id), updateByte)
		if putErr != nil {
			return putErr
		}
	}

	// 提交事务
	err = tx.Commit()
	return err
}

func Save(entity *entity.Note) error {
	entity.Id = uuid.New()
	buf, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return dao.SaveByTableNameAndKey(bean.GetDbBean(), tableName, entity.Id, buf, entityName)
}

func RenameNote(updateEntity entity.Note) error {
	if updateEntity.Id == "" {
		return errors.New("笔记ID不能为空")
	}
	dbEntity, err := FindById(updateEntity.Id)
	if err != nil {
		return err
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
	dbEntity.Alias = updateEntity.Alias
	if dbEntity.Name != updateEntity.Name {
		originalAbsPath := dbEntity.AbsPath
		// 修改了名称
		dbEntity.Path = yiuStr.GetReplaceEndStr(
			dbEntity.Path,
			string(os.PathSeparator)+dbEntity.Name,
			string(os.PathSeparator)+updateEntity.Name)
		dbEntity.AbsPath = yiuStr.GetReplaceEndStr(
			dbEntity.AbsPath,
			string(os.PathSeparator)+dbEntity.Name,
			string(os.PathSeparator)+updateEntity.Name)
		dbEntity.ParentPath = yiuStr.GetReplaceEndStr(
			dbEntity.ParentPath,
			string(os.PathSeparator)+dbEntity.Name,
			string(os.PathSeparator)+updateEntity.Name)
		renameErr := os.Rename(originalAbsPath, dbEntity.AbsPath)
		if renameErr != nil {
			return renameErr
		}
		// 如果是文件
		if dbEntity.IsDir {
			allNote, allErr := FindAll()
			if allErr != nil {
				return allErr
			}
			child := NoteUtil.GetChild(dbEntity, allNote, false)
			renameChildErr := renameChild(dbEntity, child, table)
			if renameChildErr != nil {
				return renameChildErr
			}
		}
	}

	entityMarshal, err := json.Marshal(dbEntity)
	if err != nil {
		return err
	}

	err = table.Put([]byte(dbEntity.Id), entityMarshal)
	if err != nil {
		return err
	}

	// 提交事务
	err = tx.Commit()
	return err
}

func renameChild(parent entity.Note, ch []vo.NoteTreeVo, table *bbolt.Bucket) error {
	if len(ch) == 0 {
		return nil
	}
	for i := range ch {
		item := ch[i].Data
		item.AbsPath = parent.AbsPath + string(os.PathSeparator) + item.Name
		item.Path = parent.Path + string(os.PathSeparator) + item.Name
		item.ParentAbsPath = parent.AbsPath
		itemMarshal, err := json.Marshal(item)
		if err != nil {
			return err
		}
		err = table.Put([]byte(item.Id), itemMarshal)
		if err != nil {
			return err
		}
		if len(ch[i].Child) != 0 {
			err := renameChild(item, ch[i].Child, table)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
