package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/helper"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
)

type SubDao struct {
	Db              *helper.GormHelper         `di:"~"`
	PropertyMappers *xproperty.PropertyMappers `di:"~"`
	UserDao         *UserDao                   `di:"~"`

	OrderByFunc     func(string) string `di:"-"`
	ColSubscribers  string              `di:"-"`
	ColSubscribings string              `di:"-"`
}

func NewSubDao(dic *xdi.DiContainer) *SubDao {
	repo := &SubDao{
		ColSubscribers:  "Subscribers",
		ColSubscribings: "Subscribings",
	}
	if !dic.Inject(repo) {
		log.Fatalln("Inject failed")
	}
	repo.OrderByFunc = repo.PropertyMappers.GetPropertyMapping(&dto.UserDto{}, &po.User{}).ApplyOrderBy
	return repo
}

func (s *SubDao) QuerySubscriberUsers(uid int32, page int32, limit int32, orderBy string) (users []*po.User, count int32, status database.DbStatus) {
	// https://gorm.io/docs/many_to_many.html
	user := &po.User{Uid: uid}
	if !s.UserDao.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count = int32(s.Db.Model(user).Association(s.ColSubscribers).Count()) // association pattern
	s.Db.Limit(limit).Offset((page-1)*limit).Model(user).Order(s.OrderByFunc(orderBy)).Related(&users, s.ColSubscribers)
	return users, count, database.DbSuccess
}

func (s *SubDao) QuerySubscribingUsers(uid int32, page int32, limit int32, orderBy string) (users []*po.User, count int32, status database.DbStatus) {
	user := &po.User{Uid: uid}
	if !s.UserDao.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count = int32(s.Db.Model(user).Association(s.ColSubscribings).Count())
	s.Db.Limit(limit).Offset((page-1)*limit).Model(user).Order(s.OrderByFunc(orderBy)).Related(&users, s.ColSubscribings)
	return users, count, database.DbSuccess
}

func (s *SubDao) QueryCountByUid(uid int32) (subscribingCnt int32, subscriberCnt int32, status database.DbStatus) {
	if !s.UserDao.Exist(uid) {
		return 0, 0, database.DbNotFound
	}
	user := &po.User{Uid: uid}
	subscribingCnt = int32(s.Db.Model(user).Association(s.ColSubscribings).Count())
	subscriberCnt = int32(s.Db.Model(user).Association(s.ColSubscribers).Count())
	return subscribingCnt, subscriberCnt, database.DbSuccess
}

func (s *SubDao) SubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.UserDao.Exist(toUid) || !s.UserDao.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.Db.Model(&po.User{Uid: toUid}).Association(s.ColSubscribers).Append(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (s *SubDao) UnSubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.UserDao.Exist(toUid) || !s.UserDao.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.Db.Model(&po.User{Uid: toUid}).Association(s.ColSubscribers).Delete(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}
