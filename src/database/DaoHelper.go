package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgorm"
	"github.com/jinzhu/gorm"
)

type DbHelper struct {
	*gorm.DB
}

func NewDbHelper(db *gorm.DB) *DbHelper {
	return &DbHelper {db}
}

func (db *DbHelper) QueryHelper(model interface{}, where interface{}) interface{} {
	rdb := db.Model(model).Where(where).First(where)
	if rdb.RecordNotFound() {
		return nil
	}
	return where
}

func (db *DbHelper) CountHelper(model interface{}, where interface{}) int32 {
	cnt := 0
	db.Model(model).Where(where).Count(&cnt)
	return int32(cnt)
}

func (db *DbHelper) QueryMultiHelper(model interface{}, pageSize int32, currentPage int32, where interface{}, order string, out interface{}) int32 {
	var total int32 = 0
	db.Model(model).Count(&total)
	db.Model(model).Limit(pageSize).Offset((currentPage - 1) * pageSize).Where(where).Order(order).Find(out)
	return total
}

func (db *DbHelper) SearchHelper(model interface{}, pageSize int32, currentPage int32, columns string, against string, out interface{}) int32 {
	var total int32
	rdb := db.Model(model).Where(fmt.Sprintf("MATCH (%s) AGAINST (? IN BOOLEAN MODE)", columns), against)
	rdb.Count(&total)
	rdb.Limit(pageSize).Offset((currentPage - 1) * pageSize).Find(out)
	return total
}

func (db *DbHelper) ExistHelper(model interface{}, where interface{}) bool {
	cnt := 0
	db.Model(model).Where(where).Count(&cnt)
	return cnt > 0
}

func (db *DbHelper) InsertHelper(model interface{}, object interface{}) DbStatus {
	rdb := db.Model(model).Create(object)
	if xgorm.IsMySqlDuplicateError(rdb.Error) {
		return DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return DbFailed
	}
	return DbSuccess
}

func (db *DbHelper) UpdateHelper(model interface{}, object interface{}) DbStatus {
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

func (db *DbHelper) DeleteHelper(model interface{}, object interface{}) DbStatus {
	rdb := db.Model(model).Delete(object)
	if rdb.Error != nil {
		return DbFailed
	} else if rdb.RowsAffected == 0 {
		return DbNotFound
	}
	return DbSuccess
}
