package database

import (
	. "vid/models"
)

type UserDao struct{}

const col_id string = "id"
const col_username string = "username"

func (u *UserDao) QueryAllUsers() (users []User) {
	DB.Find(&users)
	return users
}

// @return `*user` / `nil`
func (u *UserDao) QueryUser(id int) *User {
	var user User
	DB.Where(col_id+" = ?", id).Find(&user)
	if !user.CheckValid() { // PK is null (auto-increment)
		return nil
	} else {
		return &user
	}
}

// @return `*user` / `nil`
func (u *UserDao) InsertUser(user User) *User {
	if u.QueryUser(user.ID) != nil {
		return nil
	}
	DB.Create(&user)
	return u.QueryUser(user.ID)
}

// @return `*user` / `nil`
func (u *UserDao) UpdateUser(user User) *User {
	if u.QueryUser(user.ID) == nil {
		return nil
	}
	DB.Save(&user)
	return u.QueryUser(user.ID)
}

func (u *UserDao) DeleteUser(id int) bool {
	query := u.QueryUser(id)
	if query == nil {
		return false
	}
	DB.Delete(query)
	return u.QueryUser(id) == nil
}
