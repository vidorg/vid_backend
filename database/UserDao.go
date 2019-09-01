package database

import (
	. "vid/exceptions"
	. "vid/models"
)

type UserDao struct{}

const col_user_uid string = "uid"
const col_user_username string = "username"
const col_user_profile string = "profile"
const col_user_register_time string = "register_time"

// db 查询所有用户
//
// @return `[]User`
func (u *UserDao) QueryAllUsers() (users []User) {
	DB.Find(&users)
	return users
}

// db 查询 uid 用户
//
// @return `*user` `isUserExist`
func (u *UserDao) QueryUserByUid(uid int) (*User, bool) {
	var user User
	DB.Where(col_user_uid+" = ?", uid).Find(&user)
	if !user.CheckValid() { // PK is null (auto-increment)
		return nil, false
	} else {
		return &user, true
	}
}

// db 查询 username 用户
//
// @return `*user` `isUserExist`
func (u *UserDao) QueryUserByUserName(username string) (*User, bool) {
	var user User
	DB.Where(col_user_username+" = ?", username).Find(&user)
	if !user.CheckValid() { // PK is null (auto-increment)
		return nil, false
	} else {
		return &user, true
	}
}

// db 更新用户名和简介
//
// @return `*user` `err`
//
// @error `UserNotExistException` `UpdateInvalidException` `UserNameUsedException` `UpdateException` `NotUpdateException`
func (u *UserDao) UpdateUser(user User) (*User, error) {
	queryBefore, ok := u.QueryUserByUid(user.Uid)
	if !ok {
		return nil, UserNotExistException
	}
	if !user.CheckFormat() {
		return nil, UpdateInvalidException
	}
	if _, ok = u.QueryUserByUserName(user.Username); ok && user.Username != queryBefore.Username {
		return nil, UserNameUsedException
	}
	DB.Model(&user).Updates(map[string]interface{}{
		col_user_username: user.Username,
		col_user_profile:  user.Profile,
	})
	query, ok := u.QueryUserByUid(user.Uid)
	if !ok {
		return query, UpdateException
	} else {
		if queryBefore.Equals(query) {
			return query, NotUpdateException
		} else {
			return query, nil
		}
	}
}

// db 删除用户和用户密码 (cascade)
//
// @return `*user` `err`
//
// @error `UserNotExistException` `DeleteException`
func (u *UserDao) DeleteUser(uid int) (*User, error) {

	// var passDao = new(PassDao)

	query, ok := u.QueryUserByUid(uid)
	if !ok {
		return nil, UserNotExistException
	}

	// _, err := passDao.DeletePass(uid)
	DB.Delete(query)
	_, ok = u.QueryUserByUid(uid)

	if ok {
		return query, DeleteException
	} else {
		return query, nil
	}
}

// db `suberUip` 关注 `upUid`
//
// @return `err`
//
// @error `UserNotExistException` `SubscribeOneSelfException`
func (u *UserDao) SubscribeUser(suberUid int, upUid int) error {
	upUser, ok := u.QueryUserByUid(upUid)
	if !ok {
		return UserNotExistException
	}
	suberUser, ok := u.QueryUserByUid(suberUid)
	if !ok {
		return UserNotExistException
	}
	if upUid == suberUid {
		return SubscribeOneSelfException
	}
	DB.Model(upUser).Association("Subscribers").Append(suberUser)
	return nil
}

// db `suberUip` 取消关注 `upUid`
//
// @return `err`
//
// @error `UserNotExistException` `SubscribeOneSelfException`
func (u *UserDao) UnSubscribeUser(suberUid int, upUid int) error {
	upUser, ok := u.QueryUserByUid(upUid)
	if !ok {
		return UserNotExistException
	}
	suberUser, ok := u.QueryUserByUid(suberUid)
	if !ok {
		return UserNotExistException
	}
	if upUid == suberUid {
		return SubscribeOneSelfException
	}
	DB.Model(upUser).Association("Subscribers").Delete(suberUser)
	return nil
}

// db 查询 uid 的粉丝
//
// @return `user[]` `err`
//
// error `UserNotExistException`
func (u *UserDao) QuerySubscriberUsers(uid int) ([]User, error) {
	user, ok := u.QueryUserByUid(uid)
	if !ok {
		return nil, UserNotExistException
	}
	var users []User
	DB.Model(user).Related(&users, "Subscribers")
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`subscriber_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`user_uid` IN (5))
	return users, nil
}

// db 查询 uid 的关注
//
// @return `user[]` `err`
//
// @error `UserNotExistException`
func (u *UserDao) QuerySubscribingUsers(uid int) ([]User, error) {
	user, ok := u.QueryUserByUid(uid)
	if !ok {
		return nil, UserNotExistException
	}
	var users []User
	DB.Model(user).Related(&users, "Subscribings")
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`user_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`subscriber_uid` IN (5))
	return users, nil
}
