package dao

import (
	. "vid/app/database"
	"vid/app/model/po"
)

type subDao struct{}

var SubDao = new(subDao)

func (u *subDao) SubscribeUser(userUid int, upUid int) DbStatus {
	user := UserDao.QueryByUid(userUid)
	if user == nil {
		return DbNotFound
	}
	upUser := UserDao.QueryByUid(upUid)
	if upUser == nil {
		return DbNotFound
	}
	if userUid == upUid {
		return DbExtra
	}
	DB.Model(upUser).Association("Subscribers").Append(user)
	return DbSuccess
}

func (u *subDao) UnSubscribeUser(userUid int, upUid int) DbStatus {
	user := UserDao.QueryByUid(userUid)
	if user == nil {
		return DbNotFound
	}
	upUser := UserDao.QueryByUid(upUid)
	if upUser == nil {
		return DbNotFound
	}
	if userUid == upUid {
		return DbExtra
	}
	DB.Model(upUser).Association("Subscribers").Delete(user)
	return DbSuccess
}

func (u *subDao) QuerySubscriberUsers(uid int) ([]po.User, DbStatus) {
	user := UserDao.QueryByUid(uid)
	if user == nil {
		return nil, DbNotFound
	}

	var users []po.User
	DB.Model(user).Related(&users, "Subscribers")
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`subscriber_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`user_uid` IN (5))
	return users, DbSuccess
}

func (u *subDao) QuerySubscribingUsers(uid int) ([]po.User, DbStatus) {
	user := UserDao.QueryByUid(uid)
	if user == nil {
		return nil, DbNotFound
	}

	var users []po.User
	DB.Model(user).Related(&users, "Subscribings")
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`user_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`subscriber_uid` IN (5))
	return users, DbSuccess
}

func (u *subDao) QuerySubCnt(uid int) (int, int, DbStatus) {
	user := UserDao.QueryByUid(uid)
	if user == nil {
		return 0, 0, DbNotFound
	}
	subscribingCnt := DB.Model(user).Association("Subscribings").Count()
	subscriberCnt := DB.Model(user).Association("Subscribers").Count()
	return subscribingCnt, subscriberCnt, DbSuccess
}
