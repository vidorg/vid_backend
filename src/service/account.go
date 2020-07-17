package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type AccountService struct {
	Db     *database.GormHelper `di:"~"`
	Logger *logrus.Logger       `di:"~"`
}

func NewAccountService(dic *xdi.DiContainer) *AccountService {
	repo := &AccountService{}
	dic.MustInject(repo)
	return repo
}

func (a *AccountService) QueryByUsername(username string) *po.Account {
	out := a.Db.QueryFirstHelper(&po.User{}, &po.User{Username: username})
	if out == nil {
		return nil
	}
	user := out.(*po.User)
	out = a.Db.QueryFirstHelper(&po.Account{}, &po.Account{Uid: user.Uid})
	if out == nil {
		return nil
	}
	account := out.(*po.Account)
	account.User = user
	return account
}

func (a *AccountService) Insert(pass *po.Account) database.DbStatus {
	return a.Db.InsertHelper(&po.Account{}, pass) // cascade create
}

func (a *AccountService) Update(pass *po.Account) database.DbStatus {
	return a.Db.UpdateHelper(&po.Account{}, pass)
}
