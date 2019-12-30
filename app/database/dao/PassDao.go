package dao

import (
	"log"
	. "vid/app/database"
	"vid/app/model/po"
)

type passDao struct{}

var PassDao = new(passDao)

func (p *passDao) QueryByUsername(username string) *po.Password {
	user := &po.User{Username: username}
	if DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	pass := &po.Password{Uid: user.Uid}
	if DB.Where(pass).First(pass).RecordNotFound() {
		return nil
	}
	pass.User = user
	return pass
}

func (p *passDao) Insert(pass *po.Password) DbStatus {
	if err := DB.Create(pass).Error; err != nil {
		if IsDuplicateError(err) {
			return DbExisted
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}

func (p *passDao) Update(pass *po.Password) DbStatus {
	if err := DB.Model(pass).Update(pass).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}
