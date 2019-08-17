package database

import (
	"errors"
	"fmt"
	"time"
	. "vid/models"
)

type UserDao struct{}

const col_uid string = "uid"
const col_username string = "username"

// @return `[]User`
func (u *UserDao) QueryAllUsers() (users []User) {
	DB.Find(&users)
	return users
}

// @return `*user` `isUserExist`
func (u *UserDao) QueryUser(uid int) (*User, bool) {
	var user User
	DB.Where(col_uid+" = ?", uid).Find(&user)
	if !user.CheckValid() { // PK is null (auto-increment)
		return nil, false
	} else {
		return &user, true
	}
}

// @return `*user` `err`
func (u *UserDao) InsertUser(user User) (*User, error) {
	if _, ok := u.QueryUser(user.Uid); ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d already exist", user.Uid))
	}
	user.RegisterTime = time.Now()
	DB.Create(&user)
	query, ok := u.QueryUser(user.Uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d insert failed", user.Uid))
	} else {
		return query, nil
	}
}

// @return `*user` `err`
func (u *UserDao) UpdateUser(user User) (*User, error) {
	// queryBefore, ok := u.QueryUser(user.Uid)
	_, ok := u.QueryUser(user.Uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d not exist", user.Uid))
	}
	// DB.Save(&user)
	DB.Model(&user).Updates(map[string]interface{}{
		"username": user.Username,
		"profile":  user.Profile,
	})
	query, ok := u.QueryUser(user.Uid)
	if !ok {
		return query, errors.New(fmt.Sprintf("Uid: %d update failed", user.Uid))
	} else {
		// if queryBefore.Equals(query) {
		// 	return query, errors.New(fmt.Sprintf("Uid: %d not updated", user.Uid))
		// } else {
		return query, nil
		// }
	}
}

// @return `*user` `err`
func (u *UserDao) DeleteUser(uid int) (*User, error) {
	query, ok := u.QueryUser(uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d not exist", uid))
	}
	DB.Delete(query)
	_, ok = u.QueryUser(uid)
	if !ok {
		return query, errors.New(fmt.Sprintf("Uid: %d delete failed", uid))
	} else {
		return query, nil
	}
}

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

func (u *UserDao) QuerySubscriberUsers(uid int) ([]User, error) {
	user, ok := u.QueryUser(uid)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Uid: %d not exist", uid))
	}
	var users []User
	// DB.Preload("Subscribers").Find(&users, "uid = ?", uid)
	DB.Model(user).Related(&users, "Subscribers")
	// SELECT `tbl_user`.*
	// 		FROM `tbl_user` INNER JOIN `tbl_subscribe`
	// 		ON `tbl_subscribe`.`subscriber_uid` = `tbl_user`.`uid`
	// 		WHERE (`tbl_subscribe`.`user_uid` IN (5))
	return users, nil
}

func (u *UserDao) QuerySubscribingUsers(uid int) ([]User, error) {
	user, ok := u.QueryUser(uid)
	// _, ok := u.QueryUser(uid)
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
