package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
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

func (a *AccountService) Insert(account *po.Account) xstatus.DbStatus {
	return xgorm.WithDB(a.db).Insert(&po.Account{}, account) // cascade create
}

func (a *AccountService) Update(account *po.Account) xstatus.DbStatus {
	return xgorm.WithDB(a.db).Update(&po.Account{}, &po.Account{Uid: account.Uid}, account)
}
