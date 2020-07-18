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

type SubscribeService struct {
	db          *gorm.DB
	userService *UserService

	_orderByFunc func(string) string
}

func NewSubscribeService() *SubscribeService {
	return &SubscribeService{
		db:           xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService:  xdi.GetByNameForce(sn.SUserService).(*UserService),
		_orderByFunc: xproperty.GetMapperDefault(&dto.UserDto{}, &po.User{}).ApplyOrderBy,
	}
}

func (s *SubscribeService) QuerySubscriberUsers(uid int32, pageOrder *param.PageOrderParam) ([]*po.User, int32, database.DbStatus) {
	// https://gorm.io/docs/many_to_many.html
	user := &po.User{Uid: uid}
	if !s.userService.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count := int32(s.db.Model(user).Association("Subscribers").Count()) // association pattern
	users := make([]*po.User, 0)
	s.db.PageHelper(pageOrder.Limit, pageOrder.Page).Model(user).Order(s._orderByFunc(pageOrder.Order)).Related(&users, "Subscribers")
	return users, count, database.DbSuccess
}

func (s *SubscribeService) QuerySubscribingUsers(uid int32, pageOrder *param.PageOrderParam) ([]*po.User, int32, database.DbStatus) {
	user := &po.User{Uid: uid}
	if !s.userService.Exist(uid) {
		return nil, 0, database.DbNotFound
	}
	count := int32(s.db.Model(user).Association("Subscribings").Count())
	users := make([]*po.User, 0)
	s.db.PageHelper(pageOrder.Limit, pageOrder.Page).Model(user).Order(s._orderByFunc(pageOrder.Order)).Related(&users, "Subscribings")
	return users, count, database.DbSuccess
}

func (s *SubscribeService) QueryCountByUid(uid int32) (subscribingCnt int32, subscriberCnt int32, status database.DbStatus) {
	if !s.userService.Exist(uid) {
		return 0, 0, database.DbNotFound
	}
	user := &po.User{Uid: uid}
	subscribingCnt = int32(s.db.Model(user).Association("Subscribings").Count())
	subscriberCnt = int32(s.db.Model(user).Association("Subscribers").Count())
	return subscribingCnt, subscriberCnt, database.DbSuccess
}

func (s *SubscribeService) SubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.userService.Exist(toUid) || !s.userService.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.db.Model(&po.User{Uid: toUid}).Association("Subscribers").Append(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (s *SubscribeService) UnSubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.userService.Exist(toUid) || !s.userService.Exist(meUid) {
		return database.DbNotFound
	}
	ass := s.db.Model(&po.User{Uid: toUid}).Association("Subscribers").Delete(&po.User{Uid: meUid})
	if ass.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}
