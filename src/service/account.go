package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
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

func (a *AccountService) QueryByUser(user *po.User) (*po.Account, error) {
	account := &po.Account{Uid: user.Uid}
	rdb := a.db.Model(&po.Account{}).Where(&po.Account{Uid: user.Uid}).First(account)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	account.User = user
	return account, nil
}

func (a *AccountService) QueryByEmail(email string) (*po.Account, error) {
	user := &po.User{}
	rdb := a.db.Model(&po.User{}).Where(&po.User{Email: email}).First(user)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return a.QueryByUser(user)
}

func (a *AccountService) QueryByUsername(username string) (*po.Account, error) {
	user := &po.User{}
	rdb := a.db.Model(&po.User{}).Where(&po.User{Username: username}).First(user)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return a.QueryByUser(user)
}

func (a *AccountService) QueryByUid(uid uint64) (*po.Account, error) {
	user := &po.User{}
	rdb := a.db.Model(&po.User{}).Where(&po.User{Uid: uid}).First(user)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return a.QueryByUser(user)
}

func (a *AccountService) Insert(username, email, encrypted string) (xstatus.DbStatus, error) {
	account := &po.Account{
		Password: encrypted,
		User: &po.User{
			Email:    email,
			Username: username,
			Nickname: username,
		},
	}
	rdb := a.db.Model(&po.Account{}).Create(account)
	return xgorm.CreateErr(rdb)
}

func (a *AccountService) UpdatePassword(uid uint64, encrypted string) (xstatus.DbStatus, error) {
	rdb := a.db.Model(&po.Account{}).Where(&po.Account{Uid: uid}).Update("password", encrypted)
	return xgorm.UpdateErr(rdb)
}
