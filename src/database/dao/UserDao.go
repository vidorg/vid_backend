package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type UserDao struct {
	Config *config.ServerConfig `di:"~"`
	Db     *gorm.DB             `di:"~"`

	PageSize int32 `di:"-"`
}

func NewUserDao(dic *xdi.DiContainer) *UserDao {
	repo := &UserDao{}
	if !dic.Inject(repo) {
		panic("Inject failed")
	}
	repo.PageSize = repo.Config.MySqlConfig.PageSize
	return repo
}

func (u *UserDao) QueryAll(page int32) ([]*po.User, int32) {
	users := make([]*po.User, 0)
	total := PageHelper(u.Db, &po.User{}, u.PageSize, page, &po.User{}, &users)
	return users, total
}

func (u *UserDao) QueryByUid(uid int32) *po.User {
	return QueryHelper(u.Db, &po.User{}, &po.User{Uid: uid}).(*po.User)
}

func (u *UserDao) Exist(uid int32) bool {
	return ExistHelper(u.Db, &po.User{}, &po.User{Uid: uid})
}

func (u *UserDao) Update(user *po.User) database.DbStatus {
	return UpdateHelper(u.Db, &po.User{}, user)
}

func (u *UserDao) Delete(uid int32) database.DbStatus {
	ret := DeleteHelper(u.Db, &po.User{}, &po.User{Uid: uid})
	DeleteHelper(u.Db, &po.Account{}, &po.Account{Uid: uid})
	return ret
}
