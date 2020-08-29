package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/jinzhu/gorm"
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
		_orderByFunc: xgorm.OrderByFunc(xproperty.GetDefaultMapper(&dto.UserDto{}, &po.User{}).GetDict()),
	}
}

func (s *SubscribeService) getSubscriberAsso(uid int32) *gorm.Association {
	return s.db.Model(&po.User{Uid: uid}).Association("Subscribers")
}

func (s *SubscribeService) getSubscribingAsso(uid int32) *gorm.Association {
	return s.db.Model(&po.User{Uid: uid}).Association("Subscribings")
}

func (s *SubscribeService) QuerySubscriberUsers(uid int32, pp *param.PageOrderParam) (users []*po.User, total int32, status xstatus.DbStatus) {
	if !s.userService.Exist(uid) {
		return nil, 0, xstatus.DbNotFound
	}

	total = int32(s.getSubscriberAsso(uid).Count()) // association pattern
	users = make([]*po.User, 0)
	xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Model(&po.User{Uid: uid}).Order(s._orderByFunc(pp.Order)).Related(&users, "Subscribers")

	return users, total, xstatus.DbSuccess
}

func (s *SubscribeService) QuerySubscribingUsers(uid int32, pp *param.PageOrderParam) (users []*po.User, total int32, status xstatus.DbStatus) {
	if !s.userService.Exist(uid) {
		return nil, 0, xstatus.DbNotFound
	}

	total = int32(s.getSubscribingAsso(uid).Count()) // association pattern
	users = make([]*po.User, 0)
	xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Model(&po.User{Uid: uid}).Order(s._orderByFunc(pp.Order)).Related(&users, "Subscribings")

	return users, total, xstatus.DbSuccess
}

func (s *SubscribeService) QueryCountByUid(uid int32) (subscribingCnt int32, subscriberCnt int32, status xstatus.DbStatus) {
	if !s.userService.Exist(uid) {
		return 0, 0, xstatus.DbNotFound
	}

	subscribingCnt = int32(s.getSubscribingAsso(uid).Count())
	subscriberCnt = int32(s.getSubscriberAsso(uid).Count())

	return subscribingCnt, subscriberCnt, xstatus.DbSuccess
}

func (s *SubscribeService) SubscribeUser(meUid int32, toUid int32) xstatus.DbStatus {
	if !s.userService.Exist(toUid) || !s.userService.Exist(meUid) {
		return xstatus.DbNotFound
	}

	asc := s.getSubscriberAsso(toUid).Append(&po.User{Uid: meUid})
	if asc.Error != nil {
		return xstatus.DbFailed
	}
	return xstatus.DbSuccess
}

func (s *SubscribeService) UnSubscribeUser(meUid int32, toUid int32) xstatus.DbStatus {
	if !s.userService.Exist(toUid) || !s.userService.Exist(meUid) {
		return xstatus.DbNotFound
	}

	asc := s.getSubscriberAsso(toUid).Delete(&po.User{Uid: meUid})
	if asc.Error != nil {
		return xstatus.DbFailed
	}
	return xstatus.DbSuccess
}
