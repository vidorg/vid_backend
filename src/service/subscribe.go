package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
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
	orderBy     func(string) string
}

func NewSubscribeService() *SubscribeService {
	return &SubscribeService{
		db:          xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService: xdi.GetByNameForce(sn.SUserService).(*UserService),
		orderBy:     xgorm.OrderByFunc(dto.BuildUserPropertyMapper()),
	}
}

func (s *SubscribeService) subscriberAsso(uid uint64) *gorm.Association {
	return s.db.Model(&po.User{Uid: uid}).Association("Subscribers")
}

func (s *SubscribeService) subscribingAsso(uid uint64) *gorm.Association {
	return s.db.Model(&po.User{Uid: uid}).Association("Subscribings")
}

func (s *SubscribeService) QuerySubscribers(uid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.subscriberAsso(uid).Count()) // association pattern
	users := make([]*po.User, 0)
	// TODO https://gorm.io/docs/associations.html#Find-Associations
	rdb := xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Model(&po.User{Uid: uid}).Order(s.orderBy(pp.Order)).Related(&users, "Subscribers")
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return users, total, nil
}

func (s *SubscribeService) QuerySubscribings(uid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.subscribingAsso(uid).Count()) // association pattern
	users := make([]*po.User, 0)
	rdb := xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Model(&po.User{Uid: uid}).Order(s.orderBy(pp.Order)).Related(&users, "Subscribings")
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	return users, total, nil
}

func (s *SubscribeService) QueryCountByUid(uid uint64) (subscribingCnt int32, subscriberCnt int32, err error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return 0, 0, err
	} else if !ok {
		return -1, -1, nil
	}

	subscribingCnt = int32(s.subscribingAsso(uid).Count())
	subscriberCnt = int32(s.subscriberAsso(uid).Count())

	return subscribingCnt, subscriberCnt, nil
}

func (s *SubscribeService) InsertSubscribe(uid uint64, to uint64) (xstatus.DbStatus, error) {
	// TODO
	ok1, err1 := s.userService.Existed(uid)
	ok2, err2 := s.userService.Existed(to)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 || !ok2 {
		return xstatus.DbNotFound, nil
	}

	ras := s.subscriberAsso(to).Append(&po.User{Uid: uid})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
	}

	return xstatus.DbSuccess, nil
}

func (s *SubscribeService) DeleteSubscribe(uid uint64, to uint64) (xstatus.DbStatus, error) {
	// TODO
	ok1, err1 := s.userService.Existed(uid)
	ok2, err2 := s.userService.Existed(to)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 || !ok2 {
		return xstatus.DbNotFound, nil
	}

	ras := s.subscriberAsso(to).Delete(&po.User{Uid: uid})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
	}

	return xstatus.DbSuccess, nil
}
