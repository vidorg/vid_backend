package database

import (
	"errors"
	"fmt"
	"time"
	. "vid/models"
)

type PassDao struct{}

const col_pass_uid string = "uid"
const col_pass_hashpass string = "hash_pass"

// db 内部使用 查询密码项
//
// @return `isExist`
func (p *PassDao) queryPassRecord(uid int) (*Passrecord, bool) {
	var pass Passrecord
	DB.Where(col_pass_uid+" = ?", uid).Find(&pass)
	if !pass.CheckValid() {
		return nil, false
	} else {
		return &pass, true
	}
}

// db 注册 插入密码项
//
// @return `*user` `err`
//
// @error `Username: %s already exist` `Uid: %d insert failed`
func (p *PassDao) InsertUserPassRecord(username string, hashPass string) (*User, error) {

	var userDao = new(UserDao)

	if _, ok := userDao.QueryUserName(username); ok {
		return nil, errors.New(fmt.Sprintf("Username: %s already exist", username))
	}

	tx := DB.Begin()
	DB.Create(&User{
		Username:     username,
		RegisterTime: time.Now(),
	})
	queryUser, ok := userDao.QueryUserName(username)
	if !ok {
		DB.Rollback()
		return nil, errors.New(fmt.Sprintf("Uid: %d insert failed", queryUser.Uid))
	}
	DB.Create(&Passrecord{
		Uid:      queryUser.Uid,
		HashPass: hashPass,
	})

	queryPass, ok := p.queryPassRecord(queryUser.Uid)
	if !ok {
		tx.Rollback()
		return nil, errors.New(fmt.Sprintf("Uid: %d insert failed", queryPass.Uid))
	} else {
		tx.Commit()
		return queryUser, nil
	}
}

// db 登录 查询密码项
//
// @return `uid` `err`
//
// @error `Uid: %d already exist` `Uid: %d insert failed`
func (p *PassDao) QueryUidPass(uid int, hashPass string) error {
	query, ok := p.queryPassRecord(uid)
	if !ok {
		return errors.New(fmt.Sprintf("Uid: %d not exist", uid))
	}
	if query.HashPass == hashPass {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Uid: %d password error", uid))
	}
}

// db 登录 修改密码
//
// @return `uid` `err`
//
// @error `Uid: %d already exist` `Uid: %d update failed` `Uid: %d insert failed`
func (p *PassDao) UpdatePass(pass Passrecord) (int, error) {
	queryBefore, ok := p.queryPassRecord(pass.Uid)
	if !ok {
		return -1, errors.New(fmt.Sprintf("Uid: %d not exist", pass.Uid))
	}
	DB.Model(&pass).Updates(map[string]interface{}{
		col_pass_hashpass: pass.HashPass,
	})
	query, ok := p.queryPassRecord(pass.Uid)
	if !ok {
		return -1, errors.New(fmt.Sprintf("Uid: %d update failed", pass.Uid))
	} else {
		if queryBefore.HashPass == query.HashPass {
			return -1, errors.New(fmt.Sprintf("Uid: %d not updated", pass.Uid))
		} else {
			return pass.Uid, nil
		}
	}
}

// db 登录 删除密码
//
// @return `uid` `err`
//
// @error `Uid: %d already exist` `Uid: %d update failed` `Uid: %d insert failed`
func (p *PassDao) DeletePass(uid int) (int, error) {
	query, ok := p.queryPassRecord(uid)
	if !ok {
		return -1, errors.New(fmt.Sprintf("Uid: %d not exist", uid))
	}
	DB.Delete(query)
	_, ok = p.queryPassRecord(uid)
	if ok {
		return -1, errors.New(fmt.Sprintf("Uid: %d delete failed", uid))
	} else {
		return uid, nil
	}
}
