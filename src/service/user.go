package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/helper"
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
		_orderByFunc: xproperty.GetMapperDefault(&dto.UserDto{}, &po.User{}).ApplyOrderBy,
	}
}

func (u *UserService) QueryAll(pageOrder *param.PageOrderParam) (users []*po.User, total int32) {
	total = 0
	u.db.Model(&po.User{}).Count(&total)

	users = make([]*po.User, 0)
	helper.GormPager(u.db, pageOrder.Limit, pageOrder.Page).
		Model(&po.User{}).
		Order(u._orderByFunc(pageOrder.Order)).
		Find(&users)

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
	return helper.GormExist(u.db, &po.User{}, &po.User{Uid: uid})
}

func (u *UserService) Update(user *po.User) database.DbStatus {
	return helper.GormUpdate(u.db, &po.User{}, user)
}

func (u *UserService) Delete(uid int32) database.DbStatus {
	tx := u.db.Begin()

	status := helper.GormDelete(tx, &po.User{}, &po.User{Uid: uid})
	if status != database.DbSuccess {
		tx.Rollback()
		return status
	}

	status = helper.GormDelete(u.db, &po.Account{}, &po.Account{Uid: uid})
	if status != database.DbSuccess {
		tx.Rollback()
		return status
	}

	tx.Commit()
	return database.DbSuccess
}
