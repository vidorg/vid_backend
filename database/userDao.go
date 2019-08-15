package database

import (
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

// @return `*user` `beforeIsExist=false` `isOk`
func (u *UserDao) InsertUser(user User) (*User, bool, bool) {
	if _, ok := u.QueryUser(user.Uid); ok {
		return nil, true, false
	}
	user.RegisterTime = time.Now()
	DB.Create(&user)
	query, ok := u.QueryUser(user.Uid)
	return query, false, ok
}

// @return `*user`, `beforeIsExist=true` `isOk`
func (u *UserDao) UpdateUser(user User) (*User, bool, bool) {
	// queryBefore, ok := u.QueryUser(user.Uid)
	_, ok := u.QueryUser(user.Uid)
	if !ok {
		return nil, false, false
	}
	// DB.Save(&user)
	DB.Model(&user).Updates(map[string]interface{}{
		"username": user.Username,
		"profile":  user.Profile,
	})
	query, ok := u.QueryUser(user.Uid)
	// if queryBefore.Equals(query) {
	if ok {
		return query, true, false
	} else {
		return query, true, ok
	}
}

// @return `*user`, `beforeIsExist=true` `isOk`
func (u *UserDao) DeleteUser(uid int) (*User, bool, bool) {
	query, ok := u.QueryUser(uid)
	if !ok {
		return nil, false, false
	}
	DB.Delete(query)
	_, ok = u.QueryUser(uid)
	return query, true, !ok
}
