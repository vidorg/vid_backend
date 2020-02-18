package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgorm"
	"github.com/jinzhu/gorm"
)

func QueryHelper(db *gorm.DB, model interface{}, where interface{}) interface{} {
	rdb := db.Model(model).Where(where).First(where)
	if rdb.RecordNotFound() {
		return nil
	}
	return where
}

func CountHelper(db *gorm.DB, model interface{}, where interface{}) int32 {
	cnt := 0
	db.Model(model).Where(where).Count(&cnt)
	return int32(cnt)
}

func PageHelper(db *gorm.DB, model interface{}, pageSize int32, currentPage int32, where interface{}, out interface{}) int32 {
	var total int32 = 0
	db.Model(model).Count(&total)
	db.Model(model).Limit(pageSize).Offset((currentPage - 1) * pageSize).Where(where).Find(out)
	return total
}

func SearchHelper(db *gorm.DB, model interface{}, pageSize int32, currentPage int32, columns string, against string, out interface{}) int32 {
	var total int32
	rdb := db.Model(model).Where(fmt.Sprintf("MATCH (%s) AGAINST (? IN BOOLEAN MODE)", columns), against)
	rdb.Count(&total)
	rdb.Limit(pageSize).Offset((currentPage - 1) * pageSize).Find(out)
	return total
}

func ExistHelper(db *gorm.DB, model interface{}, where interface{}) bool {
	cnt := 0
	db.Model(model).Where(where).Count(&cnt)
	return cnt > 0
}

func InsertHelper(db *gorm.DB, model interface{}, object interface{}) DbStatus {
	rdb := db.Model(model).Create(object)
	if xgorm.IsMySqlDuplicateError(rdb.Error) {
		return DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return DbFailed
	}
	return DbSuccess
}

func UpdateHelper(db *gorm.DB, model interface{}, object interface{}) DbStatus {
	rdb := db.Model(model).Update(object)
	if rdb.Error != nil {
		if xgorm.IsMySqlDuplicateError(rdb.Error) {
			return DbExisted
		} else {
			return DbFailed
		}
	} else if rdb.RowsAffected == 0 {
		return DbNotFound
	}
	return DbSuccess
}

func DeleteHelper(db *gorm.DB, model interface{}, object interface{}) DbStatus {
	rdb := db.Model(model).Delete(object)
	if rdb.Error != nil {
		return DbFailed
	} else if rdb.RowsAffected == 0 {
		return DbNotFound
	}
	return DbSuccess
}
