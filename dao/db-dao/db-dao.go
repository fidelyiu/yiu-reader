package DbDao

import (
	"errors"
	"go.etcd.io/bbolt"
	"strings"
	"yiu/yiu-reader/bean"
	"yiu/yiu-reader/dao"
	"yiu/yiu-reader/model/dto"
	"yiu/yiu-reader/model/vo"
)

func FindBySearchDto(searchDto dto.DbSearchDto) (vo.DbSearchVo, error) {
	var result = vo.DbSearchVo{}
	db := bean.GetDbBean()

	if searchDto.PageSize <= 0 {
		searchDto.PageSize = 10
	}

	start := (searchDto.Page - 1) * searchDto.PageSize
	if start <= 0 {
		start = 1
	}

	end := start + searchDto.PageSize - 1
	if end <= 0 || end < start {
		end = 10
	}

	err := db.View(func(tx *bbolt.Tx) error {
		table := dao.GetTableByName(tx, searchDto.Db)
		if table == nil {
			return errors.New("不存在" + searchDto.Db + "数据库")
		}
		index := 0
		err := table.ForEach(func(k, v []byte) error {
			index++
			if searchDto.Key != "" && searchDto.Key != string(k) {
				return nil
			}
			if len(searchDto.Str) != 0 {
				for i := range searchDto.Str {
					if !strings.Contains(string(v), searchDto.Str[i]) {
						return nil
					}
				}
			}
			result.TotalRow++
			if start <= index && index <= end {
				result.Data = append(result.Data, vo.DbSearchItemVo{
					Key:   string(k),
					Value: string(v),
				})
			}
			return nil
		})
		return err
	})
	if err != nil {
		return result, err
	}

	result.TotalPage = result.TotalRow / searchDto.PageSize
	if result.TotalRow%searchDto.PageSize != 0 {
		result.TotalPage++
	}
	result.CurrentPage = searchDto.Page
	result.PageSize = searchDto.PageSize
	if 0 < result.CurrentPage && result.CurrentPage <= result.TotalPage {
		if result.CurrentPage != 1 {
			result.HasPrevious = true
		}
		if result.CurrentPage != result.TotalPage {
			result.HasNext = true
		}
	}
	return result, nil
}
