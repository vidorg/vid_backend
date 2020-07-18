package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
)

type UserService struct {
	Db      *database.GormHelper       `di:"~"`
	Logger  *logrus.Logger             `di:"~"`
	Mappers *xproperty.PropertyMappers `di:"~"`

	OrderByFunc func(string) string `di:"-"`
}

func NewUserService(dic *xdi.DiContainer) *UserService {
	repo := &UserService{}
	dic.MustInject(repo)
	repo.OrderByFunc = xproperty.GetMapperDefault(&dto.UserDto{}, &po.User{}).ApplyOrderBy
	return repo
}

func (u *UserService) QueryAll(pageOrder *param.PageOrderParam) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	total := u.Db.QueryMultiHelper(&po.User{}, pageOrder.Limit, pageOrder.Page, &po.User{}, u.OrderByFunc(pageOrder.Order), &users)
	return users, total
}

func (u *UserService) QueryByUid(uid int32) *po.User {
	out := u.Db.QueryFirstHelper(&po.User{}, &po.User{Uid: uid})
	if out == nil {
		return nil
	}
	return out.(*po.User)
}

func (u *UserService) Exist(uid int32) bool {
	return u.Db.ExistHelper(&po.User{}, &po.User{Uid: uid})
}

func (u *UserService) Update(user *po.User) database.DbStatus {
	return u.Db.UpdateHelper(&po.User{}, user)
}

func (u *UserService) Delete(uid int32) database.DbStatus {
	ret := u.Db.DeleteHelper(&po.User{}, &po.User{Uid: uid})
	if ret == database.DbSuccess {
		u.Db.DeleteHelper(&po.Account{}, &po.Account{Uid: uid})
	}
	return ret
}
