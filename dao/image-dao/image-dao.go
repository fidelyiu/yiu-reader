package ImageDao

import (
	"encoding/json"
	"errors"
	"github.com/go-basic/uuid"
	"go.etcd.io/bbolt"
	"os"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/entity"
	"yiu/yiu-reader/model/enum"
	FieldUtil "yiu/yiu-reader/util/field-util"
)

const (
	tableName  = FieldUtil.ImageCacheTable
	entityName = "图片"
)

func Save(entity *entity.Image) error {
	entity.Id = uuid.New()
	buf, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return dao.SaveByTableNameAndKey(bean.GetDbBean(), tableName, entity.Id, buf, entityName)
}

func DeleteById(id string) error {
	return dao.DeleteByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
}

func FindById(id string) (entity.Image, error) {
	v, err := dao.FindByTableNameAndKey(bean.GetDbBean(), tableName, id, entityName)
	if len(v) == 0 {
		return entity.Image{}, errors.New("[" + id + "]" + entityName + "无效")
	}
	if err != nil {
		return entity.Image{}, err
	}
	var result entity.Image
	err = json.Unmarshal([]byte(v), &result)
	return result, err
}

func DeleteFileById(id string) error {
	image, err := FindById(id)
	if err != nil {
		return err
	}
	err = DeleteFile(image)
	return err
}

func DeleteFile(image entity.Image) error {
	err := image.CheckPath()
	if err == nil {
		if image.Status == enum.ObjStatusValid {
			err = os.Remove(FieldUtil.ImageAdd + image.Path)
			if err != nil {
				return err
			}
		}
	}
	err = DeleteById(image.Id)
	return err
}

func FindByPath(path string) (entity.Image, error) {
	result := entity.Image{}
	hasData := false
	err := bean.GetDbBean().View(func(tx *bbolt.Tx) error {
		table := dao.GetTableByName(tx, tableName)
		err := table.ForEach(func(_, i []byte) error {
			var tempItem entity.Image
			err := json.Unmarshal(i, &tempItem)
			if err != nil {
				return err
			}
			if tempItem.Path == path {
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
		return result, errors.New("未找到图片")
	}
}
