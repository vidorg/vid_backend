package helper

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
)

func GormPager(db *gorm.DB, limit int32, page int32) *gorm.DB {
	return db.Limit(limit).Offset((page - 1) * limit)
}

func GormExist(db *gorm.DB, model interface{}, where interface{}) bool {
	cnt := 0
	db.Model(model).Where(where).Count(&cnt)
	return cnt > 0
}

func GormInsert(db *gorm.DB, model interface{}, object interface{}) database.DbStatus {
	rdb := db.Model(model).Create(object)

	if xgorm.IsMySqlDuplicateEntryError(rdb.Error) {
		return database.DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return database.DbFailed
	}
	return database.DbSuccess
}

func GormUpdate(db *gorm.DB, model interface{}, object interface{}) database.DbStatus {
	rdb := db.Model(model).Update(object)

	if rdb.Error != nil {
		if xgorm.IsMySqlDuplicateEntryError(rdb.Error) {
			return database.DbExisted
		} else {
			return database.DbFailed
		}
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}

func GormDelete(db *gorm.DB, model interface{}, object interface{}) database.DbStatus {
	rdb := db.Model(model).Delete(object)

	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}
