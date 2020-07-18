package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
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
	out := a.db.QueryFirstHelper(&po.User{}, &po.User{Username: username})
	if out == nil {
		return nil
	}
	user := out.(*po.User)
	out = a.db.QueryFirstHelper(&po.Account{}, &po.Account{Uid: user.Uid})
	if out == nil {
		return nil
	}
	account := out.(*po.Account)
	account.User = user
	return account
}

func (a *AccountService) Insert(pass *po.Account) database.DbStatus {
	return a.db.InsertHelper(&po.Account{}, pass) // cascade create
}

func (a *AccountService) Update(pass *po.Account) database.DbStatus {
	return a.db.UpdateHelper(&po.Account{}, pass)
}
