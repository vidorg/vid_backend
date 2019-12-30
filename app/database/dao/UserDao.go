package dao

import (
	"log"
	. "vid/app/database"
	"vid/app/model/po"
)

type userDao struct{}

var UserDao = new(userDao)

func (u *userDao) QueryAll(page int) (users []*po.User, count int) {
	DB.Model(&po.User{}).Count(&count)
	DB.Limit(PageSize).Offset((page - 1) * PageSize).Find(&users)
	return users, count
}

func (u *userDao) QueryByUid(uid int) *po.User {
	user := &po.User{Uid: uid}
	// NewRecord: 防止以主键查询时主键为 0
	if DB.NewRecord(user) || DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	return user
}

func (u *userDao) QueryByUsername(username string) *po.User {
	user := &po.User{Username: username}
	if DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	return user
}

func (u *userDao) Update(user *po.User) DbStatus {
	if DB.NewRecord(user) {
		return DbNotFound
	}
	if err := DB.Model(user).Update(user).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else if IsDuplicateError(err) {
			return DbExisted
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}

func (u *userDao) Delete(uid int) DbStatus {
	user := &po.User{Uid: uid}
	if DB.NewRecord(user) {
		return DbNotFound
	}
	if err := DB.Model(user).Update(user).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}
