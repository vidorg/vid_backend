package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/vidorg/vid_backend/lib/xgorm"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"gorm.io/gorm"
	"strings"
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
	rdb := a.db.Model(&po.Account{}).Where("uid = ?", user.Uid).First(account)
	if rdb.RowsAffected == 0 {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	account.User = user
	return account, nil
}

func (a *AccountService) QueryByEmail(email string) (*po.Account, error) {
	user := &po.User{}
	rdb := a.db.Model(&po.User{}).Where("email = ?", email).First(user)
	if rdb.RowsAffected == 0 {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return a.QueryByUser(user)
}

func (a *AccountService) QueryByUsername(username string) (*po.Account, error) {
	user := &po.User{}
	rdb := a.db.Model(&po.User{}).Where("username = ?", username).First(user)
	if rdb.RowsAffected == 0 {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return a.QueryByUser(user)
}

func (a *AccountService) QueryByUid(uid uint64) (*po.Account, error) {
	user := &po.User{}
	rdb := a.db.Model(&po.User{}).Where("uid = ?", uid).First(user)
	if rdb.RowsAffected == 0 {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	return a.QueryByUser(user)
}

func (a *AccountService) Insert(username, email, encrypted string) (xstatus.DbStatus, error) {
	tx := a.db.Begin()

	user := &po.User{
		Email:    email,
		Username: username,
		Nickname: username,
	}
	rdb := tx.Model(&po.User{}).Create(user)
	status, err := xgorm.CreateErr(rdb)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		if status == xstatus.DbExisted && err != nil {
			if strings.Contains(err.Error(), "uk_username") {
				return xstatus.DbTagA, err // username duplicated
			}
			return xstatus.DbTagB, err // email duplicated
		}
		return status, err
	}

	account := &po.Account{
		Uid:      user.Uid,
		Password: encrypted,
	}
	rdb = tx.Model(&po.Account{}).Create(account)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		return status, err
	}

	rdb = tx.Commit()
	if rdb.Error != nil {
		return status, rdb.Error
	}
	return xstatus.DbSuccess, nil
}

func (a *AccountService) UpdatePassword(uid uint64, encrypted string) (xstatus.DbStatus, error) {
	rdb := a.db.Model(&po.Account{}).Where("uid = ?", uid).Update("password", encrypted)
	return xgorm.UpdateErr(rdb)
}
