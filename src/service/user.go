package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
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

func (u *UserService) QueryAll(pageOrder *param.PageOrderParam) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	total := u.db.QueryMultiHelper(&po.User{}, pageOrder.Limit, pageOrder.Page, &po.User{}, u._orderByFunc(pageOrder.Order), &users)
	return users, total
}

func (u *UserService) QueryByUid(uid int32) *po.User {
	out := u.db.QueryFirstHelper(&po.User{}, &po.User{Uid: uid})
	if out == nil {
		return nil
	}
	return out.(*po.User)
}

func (u *UserService) Exist(uid int32) bool {
	return u.db.ExistHelper(&po.User{}, &po.User{Uid: uid})
}

func (u *UserService) Update(user *po.User) database.DbStatus {
	return u.db.UpdateHelper(&po.User{}, user)
}

func (u *UserService) Delete(uid int32) database.DbStatus {
	ret := u.db.DeleteHelper(&po.User{}, &po.User{Uid: uid})
	if ret == database.DbSuccess {
		u.db.DeleteHelper(&po.Account{}, &po.Account{Uid: uid})
	}
	return ret
}
