package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/helper"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
)

type SubscribeService struct {
	Db      *helper.GormHelper         `di:"~"`
	Logger  *logrus.Logger             `di:"~"`
	Mappers *xproperty.PropertyMappers `di:"~"`
	UserDao *UserService               `di:"~"`

	OrderByFunc func(string) string `di:"-"`
}

func NewSubscribeService(dic *xdi.DiContainer) *SubscribeService {
	repo := &SubscribeService{}
	dic.MustInject(repo)
	repo.OrderByFunc = repo.Mappers.GetPropertyMapping(&dto.UserDto{}, &po.User{}).ApplyOrderBy
	return repo
}

func (s *SubscribeService) QuerySubscriberUsers(uid int32, pageOrder *param.PageOrderParam) ([]*po.User, int32, database.DbStatus) {
	// https://gorm.io/docs/many_to_many.html
	user := &po.User{Uid: uid}
	if !s.UserDao.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count := int32(s.Db.Model(user).Association("Subscribers").Count()) // association pattern
	users := make([]*po.User, 0)
	s.Db.PageHelper(pageOrder.Limit, pageOrder.Page).Model(user).Order(s.OrderByFunc(pageOrder.Order)).Related(&users, "Subscribers")
	return users, count, database.DbSuccess
}

func (s *SubscribeService) QuerySubscribingUsers(uid int32, pageOrder *param.PageOrderParam) ([]*po.User, int32, database.DbStatus) {
	user := &po.User{Uid: uid}
	if !s.UserDao.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count := int32(s.Db.Model(user).Association("Subscribings").Count())
	users := make([]*po.User, 0)
	s.Db.PageHelper(pageOrder.Limit, pageOrder.Page).Model(user).Order(s.OrderByFunc(pageOrder.Order)).Related(&users, "Subscribings")
	return users, count, database.DbSuccess
}

func (s *SubscribeService) QueryCountByUid(uid int32) (subscribingCnt int32, subscriberCnt int32, status database.DbStatus) {
	if !s.UserDao.Exist(uid) {
		return 0, 0, database.DbNotFound
	}
	user := &po.User{Uid: uid}
	subscribingCnt = int32(s.Db.Model(user).Association("Subscribings").Count())
	subscriberCnt = int32(s.Db.Model(user).Association("Subscribers").Count())
	return subscribingCnt, subscriberCnt, database.DbSuccess
}

func (s *SubscribeService) SubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.UserDao.Exist(toUid) || !s.UserDao.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.Db.Model(&po.User{Uid: toUid}).Association("Subscribers").Append(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (s *SubscribeService) UnSubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.UserDao.Exist(toUid) || !s.UserDao.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.Db.Model(&po.User{Uid: toUid}).Association("Subscribers").Delete(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}
