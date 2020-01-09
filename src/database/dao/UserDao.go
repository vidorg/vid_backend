package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type UserDao struct {
	config *config.DatabaseConfig
	db     *gorm.DB
}

func UserRepository(config *config.DatabaseConfig) *UserDao {
	return &UserDao{
		config: config,
		db:     database.SetupDBConn(config),
	}
}

func (u *UserDao) QueryAll(page int) (users []*po.User, count int) {
	u.db.Model(&po.User{}).Count(&count)
	u.db.Limit(u.config.PageSize).Offset((page - 1) * u.config.PageSize).Find(&users)
	return users, count
}

func (u *UserDao) QueryByUid(uid int) *po.User {
	user := &po.User{Uid: uid}
	rdb := u.db.Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func (u *UserDao) Exist(uid int) bool {
	user := &po.User{Uid: uid}
	cnt := 0
	u.db.Where(user).Count(&cnt)
	return cnt > 0
}

func (u *UserDao) Update(user *po.User) database.DbStatus {
	rdb := u.db.Model(&po.User{}).Update(user)
	if rdb.Error != nil {
		if database.IsDuplicateError(rdb.Error) {
			return database.DbExisted
		} else {
			return database.DbFailed
		}
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}

func (u *UserDao) Delete(uid int) database.DbStatus {
	rdb := u.db.Model(&po.User{}).Delete(&po.User{Uid: uid})
	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}
