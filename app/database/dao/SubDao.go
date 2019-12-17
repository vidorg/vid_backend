package dao

import (
	"log"
	. "vid/app/database"
	"vid/app/model/po"
)

type subDao struct{}

var SubDao = new(subDao)

func (u *subDao) QuerySubscriberUsers(uid int, page int) ([]po.User, int, DbStatus) {
	user := &po.User{Uid: uid}
	if DB.NewRecord(user) {
		return nil, 0, DbNotFound
	}
	asDb := DB.Model(user).Association("Subscribers")
	if err := asDb.Error; err != nil {
		return nil, 0, DbNotFound
	}
	count := asDb.Count()
	var users []po.User
	DB.Limit(PageSize).Offset((page-1)*PageSize).Model(user).Related(&users, "Subscribers")
	return users, count, DbSuccess
}

func (u *subDao) QuerySubscribingUsers(uid int, page int) ([]po.User, int, DbStatus) {
	user := &po.User{Uid: uid}
	if DB.NewRecord(user) {
		return nil, 0, DbNotFound
	}
	asDb := DB.Model(user).Association("Subscribings")
	if err := asDb.Error; err != nil {
		return nil, 0, DbNotFound
	}
	count := asDb.Count()
	var users []po.User
	DB.Limit(PageSize).Offset((page-1)*PageSize).Model(user).Related(&users, "Subscribings")
	return users, count, DbSuccess
}

func (u *subDao) QuerySubCnt(uid int) (int, int, DbStatus) {
	user := &po.User{Uid: uid}
	if DB.NewRecord(user) {
		return 0, 0, DbNotFound
	}
	asDb := DB.Model(user).Association("Subscribings")
	if err := asDb.Error; err != nil {
		return 0, 0, DbNotFound
	}
	subscribingCnt := asDb.Count()
	subscriberCnt := DB.Model(user).Association("Subscribers").Count()
	return subscribingCnt, subscriberCnt, DbSuccess
}

func (u *subDao) SubscribeUser(userUid int, upUid int) DbStatus {
	user := &po.User{Uid: userUid}
	upUser := &po.User{Uid: upUid}
	if DB.NewRecord(user) || DB.NewRecord(upUser) {
		return DbNotFound
	}
	if err := DB.Model(upUser).Association("Subscribers").Append(user).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}

func (u *subDao) UnSubscribeUser(userUid int, upUid int) DbStatus {
	user := &po.User{Uid: userUid}
	upUser := &po.User{Uid: upUid}
	if DB.NewRecord(user) || DB.NewRecord(upUser) {
		return DbNotFound
	}
	if err := DB.Model(upUser).Association("Subscribers").Delete(user).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}
