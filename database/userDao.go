package database

import (
	. "vid/models"
)

type UserDao struct{}

const col_id string = "id"
const col_username string = "username"

// @return `[]User`
func (u *UserDao) QueryAllUsers() (users []User) {
	DB.Find(&users)
	return users
}

// @return `*user` `isUserExist`
func (u *UserDao) QueryUser(id int) (*User, bool) {
	var user User
	DB.Where(col_id+" = ?", id).Find(&user)
	if !user.CheckValid() { // PK is null (auto-increment)
		return nil, false
	} else {
		return &user, true
	}
}

// @return `*user` `beforeIsExist=false` `isOk`
func (u *UserDao) InsertUser(user User) (*User, bool, bool) {
	if _, ok := u.QueryUser(user.ID); ok {
		return nil, true, false
	}
	DB.Create(&user)
	query, ok := u.QueryUser(user.ID)
	return query, false, ok
}

// @return `*user`, `beforeIsExist=true` `isOk`
func (u *UserDao) UpdateUser(user User) (*User, bool, bool) {
	if _, ok := u.QueryUser(user.ID); !ok {
		return nil, false, false
	}
	DB.Save(&user)
	query, ok := u.QueryUser(user.ID)
	return query, true, ok
}

// @return `*user`, `beforeIsExist=true` `isOk`
func (u *UserDao) DeleteUser(id int) (*User, bool, bool) {
	query, ok := u.QueryUser(id)
	if !ok {
		return nil, false, false
	}
	DB.Delete(query)
	_, ok = u.QueryUser(id)
	return query, true, !ok
}
