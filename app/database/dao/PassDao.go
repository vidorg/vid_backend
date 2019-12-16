package dao

import (
	"log"
	. "vid/app/database"
	"vid/app/models/po"
)

type passDao struct{}

var PassDao = new(passDao)

func (p *passDao) Query(username string) *po.PassRecord {
	user := &po.User{Username: username}
	if DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	pass := &po.PassRecord{Uid: user.Uid}
	if DB.Where(pass).First(pass).RecordNotFound() {
		return nil
	}
	pass.User = user
	return pass
}

// DbExisted / DbFailed
func (p *passDao) Insert(pass *po.PassRecord) DbStatus {
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

// DbNotFound / DbFailed
func (p *passDao) Update(pass *po.PassRecord) DbStatus {
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
