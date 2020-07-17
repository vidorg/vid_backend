package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/jinzhu/gorm"
)

type GormHelper struct {
	*gorm.DB
}

func NewGormHelper(db *gorm.DB) *GormHelper {
	return &GormHelper{db}
}

func (db *GormHelper) PageHelper(limit int32, page int32) *GormHelper {
	return NewGormHelper(db.Limit(limit).Offset((page - 1) * limit))
}

func (db *GormHelper) QueryFirstHelper(model interface{}, where interface{}) interface{} {
	rdb := db.Model(model).Where(where).First(where)
	if rdb.RecordNotFound() {
		return nil
	}
	return where
}

func (db *GormHelper) QueryMultiHelper(model interface{}, limit int32, page int32, where interface{}, order string, out interface{}) int32 {
	var total int32 = 0
	db.Model(model).Count(&total)
	db.PageHelper(limit, page).Model(model).Where(where).Order(order).Find(out)
	return total
}

func (db *GormHelper) CountHelper(model interface{}, where interface{}) int32 {
	cnt := 0
	db.Model(model).Where(where).Count(&cnt)
	return int32(cnt)
}

func (db *GormHelper) ExistHelper(model interface{}, where interface{}) bool {
	return db.CountHelper(model, where) > 0
}

func (db *GormHelper) SearchHelper(model interface{}, limit int32, page int32, columns string, against string, out interface{}) int32 {
	var total int32
	rdb := db.Model(model).Where(fmt.Sprintf("MATCH (%s) AGAINST (? IN BOOLEAN MODE)", columns), against)
	rdb.Count(&total)
	rdb.Limit(limit).Offset((page - 1) * limit).Find(out)
	return total
}

func (db *GormHelper) InsertHelper(model interface{}, object interface{}) DbStatus {
	rdb := db.Model(model).Create(object)
	if xgorm.IsMySqlDuplicateError(rdb.Error) {
		return DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return DbFailed
	}
	return DbSuccess
}

func (db *GormHelper) UpdateHelper(model interface{}, object interface{}) DbStatus {
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

func (db *GormHelper) DeleteHelper(model interface{}, object interface{}) DbStatus {
	rdb := db.Model(model).Delete(object)
	if rdb.Error != nil {
		return DbFailed
	} else if rdb.RowsAffected == 0 {
		return DbNotFound
	}
	return DbSuccess
}

func (db *GormHelper) Begin() *GormHelper {
	return NewGormHelper(db.DB.Begin())
}

func (db *GormHelper) Commit() *GormHelper {
	return NewGormHelper(db.DB.Commit())

}
func (db *GormHelper) Rollback() *GormHelper {
	return NewGormHelper(db.DB.Rollback())
}

func (db *GormHelper) RollbackUnlessCommitted() *GormHelper {
	return NewGormHelper(db.DB.RollbackUnlessCommitted())
}
