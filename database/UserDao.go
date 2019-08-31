package database

import (
	"errors"
	"fmt"
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
func (u *UserDao) QueryUser(uid int) (*User, bool) {
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
// @error `Uid: %d not exist` `Uid: %d update failed` `Uid: %d not updated`
func (u *UserDao) UpdateUser(user User) (*User, error) {
	queryBefore, ok := u.QueryUser(user.Uid)
	// _, ok := u.QueryUser(user.Uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d not exist", user.Uid))
	}
	// DB.Save(&user)
	DB.Model(&user).Updates(map[string]interface{}{
		col_user_username: user.Username,
		col_user_profile:  user.Profile,
	})
	query, ok := u.QueryUser(user.Uid)
	if !ok {
		return query, errors.New(fmt.Sprintf("Uid: %d update failed", user.Uid))
	} else {
		if queryBefore.Equals(query) {
			return query, errors.New(fmt.Sprintf("Uid: %d not updated", user.Uid))
		} else {
			return query, nil
		}
	}
}

// db 删除用户和用户密码 (cascade)
//
// @return `*user` `err`
//
// @error `Uid: %d not exist` `Uid: %d delete failed`
func (u *UserDao) DeleteUser(uid int) (*User, error) {

	// var passDao = new(PassDao)

	query, ok := u.QueryUser(uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d not exist", uid))
	}

	// _, err := passDao.DeletePass(uid)
	DB.Delete(query)
	_, ok = u.QueryUser(uid)

	if ok {
		return query, errors.New(fmt.Sprintf("Uid: %d delete failed", uid))
	} else {
		return query, nil
	}
}

// db `suberUip` 关注 `upUid`
//
// @return `err`
//
// @error `Uid: %d not exist` `Cound not subscribe to oneself`
func (u *UserDao) SubscribeUser(upUid int, suberUid int) error {
	upUser, ok := u.QueryUser(upUid)
	if !ok {
		return errors.New(fmt.Sprintf("Uid: %d not exist", upUid))
	}
	suberUser, ok := u.QueryUser(suberUid)
	if !ok {
		return errors.New(fmt.Sprintf("Uid: %d not exist", suberUid))
	}
	if upUid == suberUid {
		return errors.New(fmt.Sprintf("Cound not subscribe to oneself"))
	}
	DB.Model(upUser).Association("Subscribers").Append(suberUser)
	return nil
}

// db `suberUip` 取消关注 `upUid`
//
// @return `err`
//
// @error `Uid: %d not exist` `Cound not subscribe to oneself`
func (u *UserDao) UnSubscribeUser(upUid int, suberUid int) error {
	upUser, ok := u.QueryUser(upUid)
	if !ok {
		return errors.New(fmt.Sprintf("Uid: %d not exist", upUid))
	}
	suberUser, ok := u.QueryUser(suberUid)
	if !ok {
		return errors.New(fmt.Sprintf("Uid: %d not exist", suberUid))
	}
	if upUid == suberUid {
		return errors.New(fmt.Sprintf("Cound not subscribe to oneself"))
	}
	DB.Model(upUser).Association("Subscribers").Delete(suberUser)
	return nil
}

// db 查询 uid 的粉丝
//
// @return `user[]` `err`
//
// error `Uid: %d not exist`
func (u *UserDao) QuerySubscriberUsers(uid int) ([]User, error) {
	user, ok := u.QueryUser(uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d not exist", uid))
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
// @error `Uid: %d not exist`
func (u *UserDao) QuerySubscribingUsers(uid int) ([]User, error) {
	user, ok := u.QueryUser(uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d not exist", uid))
	}
	var users []User
	DB.Model(user).Related(&users, "Subscribings")
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`user_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`subscriber_uid` IN (5))
	return users, nil
}