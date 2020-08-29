package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type UserService struct {
	db *gorm.DB

	_orderByFunc func(string) string
}

func NewUserService() *UserService {
	return &UserService{
		db:           xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		_orderByFunc: xgorm.OrderByFunc(xproperty.GetDefaultMapper(&dto.UserDto{}, &po.User{}).GetDict()),
	}
}

func (u *UserService) QueryAll(pp *param.PageOrderParam) (users []*po.User, total int32) {
	total = 0
	u.db.Model(&po.User{}).Count(&total)

	users = make([]*po.User, 0)
	xgorm.WithDB(u.db).Pagination(pp.Limit, pp.Page).Model(&po.User{}).Order(u._orderByFunc(pp.Order)).Find(&users)

	return users, total
}

func (u *UserService) QueryByUid(uid int32) *po.User {
	user := &po.User{}
	rdb := u.db.Model(&po.User{}).Where(&po.User{Uid: uid}).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func (u *UserService) Exist(uid int32) bool {
	cnt := 0
	u.db.Model(&po.User{}).Where(&po.User{Uid: uid}).Count(&cnt)
	return cnt > 0
}

func (u *UserService) Update(user *po.User) xstatus.DbStatus {
	rdb := u.db.Model(&po.User{}).Where(&po.User{Uid: user.Uid}).Update(user)
	status, _ := xgorm.UpdateErr(rdb)
	return status
}

func (u *UserService) Delete(uid int32) xstatus.DbStatus {
	tx := u.db.Begin()

	rdb := tx.Model(&po.User{}).Where(&po.User{Uid: uid}).Delete(&po.User{Uid: uid})
	status, _ := xgorm.DeleteErr(rdb)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		return status
	}

	rdb = tx.Model(&po.Account{}).Where(&po.Account{Uid: uid}).Delete(&po.Account{Uid: uid})
	status, _ = xgorm.DeleteErr(rdb)
	if status != xstatus.DbSuccess {
		tx.Rollback()
		return status
	}

	tx.Commit()
	return xstatus.DbSuccess
}
