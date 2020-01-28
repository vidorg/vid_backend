package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type PassDao struct {
	Db *gorm.DB `di:"~"`
}

func NewPassDao(dic *xdi.DiContainer) *PassDao {
	repo := &PassDao{}
	if !dic.Inject(repo) {
		panic("Inject failed")
	}
	return repo
}

func (p *PassDao) QueryByUsername(username string) *po.Account {
	user := &po.User{Username: username}
	rdb := p.Db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	pass := &po.Account{Uid: user.Uid}
	rdb = p.Db.Model(&po.Account{}).Where(pass).First(pass)
	if rdb.RecordNotFound() {
		return nil
	}
	pass.User = user
	return pass
}

func (p *PassDao) Insert(pass *po.Account) database.DbStatus {
	rdb := p.Db.Model(&po.Account{}).Create(pass) // cascade create
	if database.IsDuplicateError(rdb.Error) {
		return database.DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (p *PassDao) Update(pass *po.Account) database.DbStatus {
	rdb := p.Db.Model(&po.Account{}).Update(pass)
	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}
