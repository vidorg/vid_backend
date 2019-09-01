package database

import (
	"time"
	. "vid/exceptions"
	. "vid/models"
)

type PassDao struct{}

const col_pass_uid string = "uid"
const col_pass_encryptedPass string = "encrypted_pass"

// db 内部使用 查询密码项
//
// @return `isExist`
func (p *PassDao) queryPassRecord(uid int) (*PassRecord, bool) {
	var pass PassRecord
	DB.Where(col_pass_uid+" = ?", uid).Find(&pass)
	if !pass.CheckValid() {
		return nil, false
	} else {
		return &pass, true
	}
}

// db 注册 插入用户和密码项
//
// @return `*user` `err`
//
// @error `UserExistException` `InsertException`
func (p *PassDao) InsertUserPassRecord(username string, encryptedPass string) (*User, error) {

	var userDao = new(UserDao)

	if _, ok := userDao.QueryUserByUserName(username); ok {
		return nil, UserExistException
	}

	tx := DB.Begin()
	DB.Create(&User{
		Username:     username,
		RegisterTime: time.Now(),
	})
	queryUser, ok := userDao.QueryUserByUserName(username)
	if !ok {
		DB.Rollback()
		return nil, InsertException
	}
	DB.Create(&PassRecord{
		Uid:           queryUser.Uid,
		EncryptedPass: encryptedPass,
	})

	_, ok = p.queryPassRecord(queryUser.Uid)
	if !ok {
		tx.Rollback()
		return nil, InsertException
	} else {
		tx.Commit()
		return queryUser, nil
	}
}

// db 登录 查询密码项
//
// @return `*user` `*pass` `isExist`
func (p *PassDao) QueryPassRecordByUsername(username string) (*User, *PassRecord, bool) {

	var userDao = new(UserDao)
	user, ok := userDao.QueryUserByUserName(username)
	if !ok {
		return nil, nil, false
	}

	var pass PassRecord
	DB.Where(col_pass_uid+" = ?", user.Uid).Find(&pass)
	if !pass.CheckValid() {
		return nil, nil, false
	} else {
		return user, &pass, true
	}
}

// db 登录 修改密码
//
// @return `uid` `err`
//
// @error `UserExistException` `UpdateException` `NotUpdateException`
func (p *PassDao) UpdatePass(pass PassRecord) (int, error) {
	queryBefore, ok := p.queryPassRecord(pass.Uid)
	if !ok {
		return -1, UserExistException
	}
	DB.Model(&pass).Updates(map[string]interface{}{
		col_pass_encryptedPass: pass.EncryptedPass,
	})
	query, ok := p.queryPassRecord(pass.Uid)
	if !ok {
		return -1, UpdateException
	} else {
		if queryBefore.EncryptedPass == query.EncryptedPass {
			return -1, NotUpdateException
		} else {
			return pass.Uid, nil
		}
	}
}
