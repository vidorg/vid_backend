package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type UserDao struct {
	db       *gorm.DB
	pageSize int
}

func UserRepository(config *config.DatabaseConfig) *UserDao {
	return &UserDao{
		db:       database.SetupDBConn(config),
		pageSize: config.PageSize,
	}
}

func (u *UserDao) QueryAll(page int) (users []*po.User, count int) {
	u.db.Model(&po.User{}).Count(&count)
	u.db.Model(&po.User{}).Limit(u.pageSize).Offset((page - 1) * u.pageSize).Find(&users)
	return users, count
}

func (u *UserDao) QueryByUid(uid int) *po.User {
	user := &po.User{Uid: uid}
	rdb := u.db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func (u *UserDao) QueryByUsername(username string) *po.User {
	user := &po.User{Username: username}
	rdb := u.db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func (u *UserDao) Exist(uid int) bool {
	user := &po.User{Uid: uid}
	rdb := u.db.Model(&po.User{}).Where(user)
	return !rdb.RecordNotFound()
}

func (u *UserDao) Update(user *po.User) database.DbStatus {
	rdb := u.db.Model(&po.User{}).Update(user)
	if rdb.Error != nil {
		if database.IsDuplicateError(rdb.Error) {
			return database.DbExisted
		} else if database.IsNotFoundError(rdb.Error) {
			return database.DbNotFound
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}

func (u *UserDao) Delete(uid int) database.DbStatus {
	rdb := u.db.Model(&po.User{}).Delete(&po.User{Uid: uid})
	if rdb.Error != nil {
		if database.IsNotFoundError(rdb.Error) {
			return database.DbNotFound
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}
