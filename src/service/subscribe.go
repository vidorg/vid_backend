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

func (s *SubscribeService) getSubscriberAsso(uid int32) *gorm.Association {
	return s.db.Model(&po.User{Uid: uid}).Association("Subscribers")
}

func (s *SubscribeService) getSubscribingAsso(uid int32) *gorm.Association {
	return s.db.Model(&po.User{Uid: uid}).Association("Subscribings")
}

func (s *SubscribeService) QuerySubscriberUsers(uid int32, pageOrder *param.PageOrderParam) (users []*po.User, total int32, status database.DbStatus) {
	if !s.userService.Exist(uid) {
		return nil, 0, database.DbNotFound
	}

	total = int32(s.getSubscriberAsso(uid).Count()) // association pattern
	users = make([]*po.User, 0)
	helper.GormPager(s.db, pageOrder.Limit, pageOrder.Page).
		Model(&po.User{Uid: uid}).
		Order(s._orderByFunc(pageOrder.Order)).
		Related(&users, "Subscribers")

	return users, total, database.DbSuccess
}

func (s *SubscribeService) QuerySubscribingUsers(uid int32, pageOrder *param.PageOrderParam) (users []*po.User, total int32, status database.DbStatus) {
	if !s.userService.Exist(uid) {
		return nil, 0, database.DbNotFound
	}

	total = int32(s.getSubscribingAsso(uid).Count()) // association pattern
	users = make([]*po.User, 0)
	helper.GormPager(s.db, pageOrder.Limit, pageOrder.Page).
		Model(&po.User{Uid: uid}).
		Order(s._orderByFunc(pageOrder.Order)).
		Related(&users, "Subscribings")

	return users, total, database.DbSuccess
}

func (s *SubscribeService) QueryCountByUid(uid int32) (subscribingCnt int32, subscriberCnt int32, status database.DbStatus) {
	if !s.userService.Exist(uid) {
		return 0, 0, database.DbNotFound
	}

	subscribingCnt = int32(s.getSubscribingAsso(uid).Count())
	subscriberCnt = int32(s.getSubscriberAsso(uid).Count())

	return subscribingCnt, subscriberCnt, database.DbSuccess
}

func (s *SubscribeService) SubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.userService.Exist(toUid) || !s.userService.Exist(meUid) {
		return database.DbNotFound
	}

	asc := s.getSubscriberAsso(toUid).Append(&po.User{Uid: meUid})
	if asc.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (s *SubscribeService) UnSubscribeUser(meUid int32, toUid int32) database.DbStatus {
	if !s.userService.Exist(toUid) || !s.userService.Exist(meUid) {
		return database.DbNotFound
	}

	asc := s.getSubscriberAsso(toUid).Delete(&po.User{Uid: meUid})
	if asc.Error != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}
