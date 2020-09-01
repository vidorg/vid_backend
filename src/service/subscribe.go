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
	tblName     string
}

func NewSubscribeService() *SubscribeService {
	return &SubscribeService{
		db:          xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService: xdi.GetByNameForce(sn.SUserService).(*UserService),
		orderBy:     xgorm.OrderByFunc(dto.BuildUserPropertyMapper()),
		tblName:     "tbl_subscribe",
	}
}

func (s *SubscribeService) subscriberAssoDB(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Subscribers") // association pattern
}

func (s *SubscribeService) subscriberAsso(uid uint64) *gorm.Association {
	return s.subscriberAssoDB(s.db, uid)
}

func (s *SubscribeService) subscribingAssoDB(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Subscribings") // association pattern
}

func (s *SubscribeService) subscribingAsso(uid uint64) *gorm.Association {
	return s.subscribingAssoDB(s.db, uid)
}

func (s *SubscribeService) QuerySubscribers(uid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.subscriberAsso(uid).Count())
	users := make([]*po.User, 0)
	rac := s.subscriberAssoDB(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderBy(pp.Order)), uid).Find(&users)
	if rac.Error != nil {
		return nil, 0, rac.Error
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

	total := int32(s.subscribingAsso(uid).Count())
	users := make([]*po.User, 0)
	rac := s.subscribingAssoDB(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderBy(pp.Order)), uid).Find(&users)
	if rac.Error != nil {
		return nil, 0, rac.Error
	}

	return users, total, nil
}

func (s *SubscribeService) QueryCountByUids(uids []uint64) ([]*[2]int32, error) {
	if len(uids) == 0 {
		return []*[2]int32{}, nil
	}

	type result struct {
		Id  uint64
		Cnt int32
	}

	// subscribing
	sp := strings.Builder{}
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("from_uid = %d OR ", uid))
	}
	where := sp.String()[:sp.Len()-4]
	subings := make([]*result, 0)
	rdb := s.db.Table(s.tblName).Select("from_uid as id, count(*) as cnt").Where(where).Group("from_uid").Scan(&subings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	// subscriber
	sp.Reset()
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("to_uid = %d OR ", uid))
	}
	where = sp.String()[:sp.Len()-4]
	subers := make([]*result, 0)
	rdb = s.db.Table(s.tblName).Select("to_uid as id, count(*) as cnt").Where(where).Group("to_uid").Scan(&subers)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64][2]int32, len(uids))
	for _, subing := range subings {
		bucket[subing.Id] = [2]int32{subing.Cnt, 0}
	}
	for _, suber := range subers {
		if arr, ok := bucket[suber.Id]; !ok {
			bucket[suber.Id] = [2]int32{0, suber.Cnt}
		} else {
			bucket[suber.Id] = [2]int32{arr[0], suber.Cnt}
		}
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
	if len(uids) == 0 {
		return []*[2]bool{}, nil
	}

	type result struct {
		FromUid uint64
		ToUid   uint64
	}

	// subscribing
	sp := strings.Builder{}
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("to_uid = %d OR ", uid))
	}
	where := sp.String()[:sp.Len()-4]
	subings := make([]*result, 0)
	rdb := s.db.Table(s.tblName).Select("from_uid, to_uid").Where("from_uid = ?", me).Where(where).Group("from_uid, to_uid").Scan(&subings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	// subscriber
	sp.Reset()
	for _, uid := range uids {
		sp.WriteString(fmt.Sprintf("from_uid = %d OR ", uid))
	}
	where = sp.String()[:sp.Len()-4]
	subers := make([]*result, 0)
	rdb = s.db.Table(s.tblName).Select("from_uid, to_uid").Where("to_uid = ?", me).Where(where).Group("from_uid, to_uid").Scan(&subers)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64][2]bool, len(uids))
	for _, subing := range subings {
		bucket[subing.ToUid] = [2]bool{true, false}
	}
	for _, suber := range subers {
		if arr, ok := bucket[suber.FromUid]; !ok {
			bucket[suber.FromUid] = [2]bool{false, true}
		} else {
			bucket[suber.FromUid] = [2]bool{arr[0], true}
		}
	}

	out := make([]*[2]bool, len(uids))
	for idx, uid := range uids {
		arr, ok := bucket[uid]
		if ok {
			out[idx] = &arr
		} else {
			out[idx] = &[2]bool{}
		}
	}

	return out, nil
}

func (s *SubscribeService) InsertSubscribe(uid uint64, to uint64) (xstatus.DbStatus, error) {
	ok1, err1 := s.userService.Existed(uid)
	ok2, err2 := s.userService.Existed(to)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 || !ok2 {
		return xstatus.DbNotFound, nil
	}

	cnt := 0
	rdb := s.db.Table(s.tblName).Where("from_uid = ? AND to_uid = ?", uid, to).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt > 0 {
		return xstatus.DbExisted, nil
	}

	ras := s.subscriberAsso(to).Append(&po.User{Uid: uid})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
	}

	return xstatus.DbSuccess, nil
}

func (s *SubscribeService) DeleteSubscribe(uid uint64, to uint64) (xstatus.DbStatus, error) {
	ok1, err1 := s.userService.Existed(uid)
	ok2, err2 := s.userService.Existed(to)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 || !ok2 {
		return xstatus.DbNotFound, nil
	}

	cnt := 0
	rdb := s.db.Table(s.tblName).Where("from_uid = ? AND to_uid = ?", uid, to).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt == 0 {
		return xstatus.DbTagA, nil
	}

	ras := s.subscriberAsso(to).Delete(&po.User{Uid: uid})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
	}

	return xstatus.DbSuccess, nil
}
