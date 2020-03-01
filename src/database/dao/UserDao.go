package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
)

type UserDao struct {
	Config          *config.ServerConfig       `di:"~"`
	Db              *database.DbHelper         `di:"~"`
	PropertyMappers *xproperty.PropertyMappers `di:"~"`

	PageSize    int32               `di:"-"`
	OrderByFunc func(string) string `di:"-"`
}

func NewUserDao(dic *xdi.DiContainer) *UserDao {
	repo := &UserDao{}
	if !dic.Inject(repo) {
		log.Fatalln("Inject failed")
	}
	repo.PageSize = repo.Config.MySqlConfig.PageSize
	repo.OrderByFunc = repo.PropertyMappers.GetPropertyMapping(&dto.UserDto{}, &po.User{}).ApplyOrderBy
	return repo
}

func (u *UserDao) QueryAll(page int32, orderBy string) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	log.Println(u.OrderByFunc(orderBy))
	total := u.Db.QueryMultiHelper(&po.User{}, u.PageSize, page, &po.User{}, u.OrderByFunc(orderBy), &users)
	return users, total
}

func (u *UserDao) QueryByUid(uid int32) *po.User {
	out := u.Db.QueryHelper(&po.User{}, &po.User{Uid: uid})
	if out == nil {
		return nil
	}
	return out.(*po.User)
}

func (u *UserDao) Exist(uid int32) bool {
	return u.Db.ExistHelper(&po.User{}, &po.User{Uid: uid})
}

func (u *UserDao) Update(user *po.User) database.DbStatus {
	return u.Db.UpdateHelper(&po.User{}, user)
}

func (u *UserDao) Delete(uid int32) database.DbStatus {
	ret := u.Db.DeleteHelper(&po.User{}, &po.User{Uid: uid})
	if ret == database.DbSuccess {
		u.Db.DeleteHelper(&po.Account{}, &po.Account{Uid: uid})
	}
	return ret
}
