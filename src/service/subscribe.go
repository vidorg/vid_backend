package service

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"strings"
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

type groupByResult struct {
	Id  uint64
	Cnt int32
}

func (s *SubscribeService) QueryCountByUids(uids []uint64) ([]*[2]int32, error) {
	if len(uids) == 0 {
		return []*[2]int32{}, nil
	}

	// subscribing
	sp := strings.Builder{}
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("from_uid = %d OR ", uid))
	}
	where := sp.String()
	where = where[:len(where)-4]

	subings := make([]*groupByResult, 0)
	rdb := s.db.Table("tbl_subscribe").Select("from_uid as id, count(*) as cnt").Where(where).Group("from_uid").Scan(&subings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	// subscriber
	sp = strings.Builder{}
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("to_uid = %d OR ", uid))
	}
	where = sp.String()
	where = where[:len(where)-4]

	subers := make([]*groupByResult, 0)
	rdb = s.db.Table("tbl_subscribe").Select("to_uid as id, count(*) as cnt").Where(where).Group("to_uid").Scan(&subers)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64][2]int32, len(uids))
	for _, subing := range subings {
		a, ok := bucket[subing.Id]
		if !ok {
			a = [2]int32{}
		}
		a[0] = subing.Cnt
		bucket[subing.Id] = a
	}
	for _, suber := range subers {
		a, ok := bucket[suber.Id]
		if !ok {
			a = [2]int32{}
		}
		a[1] = suber.Cnt
		bucket[suber.Id] = a
	}

	out := make([]*[2]int32, len(uids))
	for idx, uid := range uids {
		arr, ok := bucket[uid]
		if ok {
			out[idx] = &arr
		}
	}

	return out, nil
}

func (s *SubscribeService) CheckSubscribeByUids(me uint64, uids []uint64) ([]*[2]bool, error) {
	// TODO
	return nil, nil
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
