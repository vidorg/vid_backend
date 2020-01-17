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

func NewUserDao(dic xdi.DiContainer) *UserDao {
	repo := &UserDao{}
	dic.Inject(repo)
	if xdi.HasNilDi(repo) {
		panic("Has nil di field")
	}

	repo.PageSize = repo.Config.MySqlConfig.PageSize
	return repo
}

func (u *UserDao) QueryAll(page int32) (users []*po.User, count int32) {
	u.Db.Model(&po.User{}).Count(&count)
	u.Db.Model(&po.User{}).Limit(u.PageSize).Offset((page - 1) * u.PageSize).Find(&users)
	return users, count
}

func (u *UserDao) QueryByUid(uid int32) *po.User {
	user := &po.User{Uid: uid}
	rdb := u.Db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	return user
}

func (u *UserDao) Exist(uid int32) bool {
	user := &po.User{Uid: uid}
	cnt := 0
	u.Db.Model(&po.User{}).Where(user).Count(&cnt)
	return cnt > 0
}

func (u *UserDao) Update(user *po.User) database.DbStatus {
	rdb := u.Db.Model(&po.User{}).Update(user)
	if rdb.Error != nil {
		if database.IsDuplicateError(rdb.Error) {
			return database.DbExisted
		} else {
			return database.DbFailed
		}
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}

func (u *UserDao) Delete(uid int32) database.DbStatus {
	rdb := u.Db.Model(&po.User{}).Delete(&po.User{Uid: uid})
	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	u.Db.Delete(&po.PassRecord{Uid: uid})
	return database.DbSuccess
}
