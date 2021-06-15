package WorkspaceDao

import (
	"encoding/json"
	"errors"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/entity"
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
