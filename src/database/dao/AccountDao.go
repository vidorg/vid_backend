package dao

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type AccountDao struct {
	Db *gorm.DB `di:"~"`
}

func NewPassDao(dic *xdi.DiContainer) *AccountDao {
	repo := &AccountDao{}
	if !dic.Inject(repo) {
		panic("Inject failed")
	}
	return repo
}

func (a *AccountDao) QueryByUsername(username string) *po.Account {
	user := &po.User{Username: username}
	rdb := a.Db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	account := &po.Account{Uid: user.Uid}
	rdb = a.Db.Model(&po.Account{}).Where(account).First(account)
	if rdb.RecordNotFound() {
		return nil
	}
	account.User = user
	return account
}

func (a *AccountDao) Insert(pass *po.Account) database.DbStatus {
	rdb := a.Db.Model(&po.Account{}).Create(pass) // cascade create
	if xgorm.IsMySqlDuplicateError(rdb.Error) {
		return database.DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (a *AccountDao) Update(pass *po.Account) database.DbStatus {
	rdb := a.Db.Model(&po.Account{}).Update(pass)
	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}
