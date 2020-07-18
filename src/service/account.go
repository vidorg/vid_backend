package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/helper"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService() *AccountService {
	return &AccountService{
		db: xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
	}
}

func (a *AccountService) QueryByUsername(username string) *po.Account {
	user := &po.User{}
	rdb := a.db.Model(&po.User{}).Where(&po.User{Username: username}).First(user)
	if rdb.RecordNotFound() {
		return nil
	}

	account := &po.Account{}
	rdb = a.db.Model(&po.Account{}).Where(&po.Account{Uid: user.Uid}).First(account)
	if rdb.RecordNotFound() {
		return nil
	}

	account.User = user
	return account
}

func (a *AccountService) Insert(pass *po.Account) database.DbStatus {
	return helper.GormInsert(a.db, &po.Account{}, pass) // cascade create
}

func (a *AccountService) Update(pass *po.Account) database.DbStatus {
	return helper.GormUpdate(a.db, &po.Account{}, pass)
}
