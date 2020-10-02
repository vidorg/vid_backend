package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/vidorg/vid_backend/lib/xgorm"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"gorm.io/gorm"
)

type SubscribeService struct {
	db             *gorm.DB
	common         *CommonService
	userService    *UserService
	orderbyService *OrderbyService
}

func NewSubscribeService() *SubscribeService {
	return &SubscribeService{
		db:             xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		common:         xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		userService:    xdi.GetByNameForce(sn.SUserService).(*UserService),
		orderbyService: xdi.GetByNameForce(sn.SOrderbyService).(*OrderbyService),
	}
}

func (s *SubscribeService) table() *gorm.DB {
	return s.db.Table("tbl_subscribe")
}

func (s *SubscribeService) subscriberAsso(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Subscribers") // association pattern
}

func (s *SubscribeService) subscribingAsso(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Subscribings") // association pattern
}

func (s *SubscribeService) QuerySubscribers(uid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.subscriberAsso(s.db, uid).Count())
	users := make([]*po.User, 0)
	err = s.subscriberAsso(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderbyService.SubscribeForUser(pp.Order)), uid).Find(&users)
	if err != nil {
		return nil, 0, err
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

	total := int32(s.subscribingAsso(s.db, uid).Count())
	users := make([]*po.User, 0)
	err = s.subscribingAsso(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderbyService.SubscribeForUser(pp.Order)), uid).Find(&users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
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

	cnt := int64(0)
	rdb := s.table().Where("from_uid = ? AND to_uid = ?", uid, to).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt > 0 {
		return xstatus.DbExisted, nil // existed
	}

	err := s.subscriberAsso(s.db, to).Append(&po.User{Uid: uid})
	if err != nil {
		return xstatus.DbFailed, err
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

	cnt := int64(0)
	rdb := s.table().Where("from_uid = ? AND to_uid = ?", uid, to).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt == 0 {
		return xstatus.DbTagA, nil // not found
	}

	err := s.subscriberAsso(s.db, to).Delete(&po.User{Uid: uid})
	if err != nil {
		return xstatus.DbFailed, err
	}

	return xstatus.DbSuccess, nil
}

func (s *SubscribeService) CheckSubscribe(me uint64, uids []uint64) (ingerPairs []*[2]bool, err error) {
	if len(uids) == 0 {
		return []*[2]bool{}, nil
	}

	subscribings := make([]*_IdScanResult, 0)
	where := s.common.BuildOrExpr("to_uid", uids)
	// TODO tx.Statement.Selects = []string{strings.Join(fields, " ")}
	rdb := s.table().Select("to_uid as id").Where("from_uid = ?", me).Where(where).Scan(&subscribings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	subscribers := make([]*_IdScanResult, 0)
	where = s.common.BuildOrExpr("from_uid", uids)
	rdb = s.table().Select("from_uid as id").Where("to_uid = ?", me).Where(where).Scan(&subscribers)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64][2]bool, len(uids))
	for _, r := range subscribings {
		bucket[r.Id] = [2]bool{true, false}
	}
	for _, r := range subscribers {
		_, ok := bucket[r.Id]
		bucket[r.Id] = [2]bool{ok, true} // use ok directly
	}

	out := make([]*[2]bool, len(uids))
	for idx, uid := range uids {
		if arr, ok := bucket[uid]; ok {
			out[idx] = &arr
		} else {
			out[idx] = &[2]bool{false, false}
		}
	}
	return out, nil
}

func (s *SubscribeService) QuerySubscribeCount(uids []uint64) (ingerPairs []*[2]int32, err error) {
	if len(uids) == 0 {
		return []*[2]int32{}, nil
	}

	subscribings := make([]*_IdCntScanResult, 0)
	where := s.common.BuildOrExpr("from_uid", uids)
	rdb := s.table().Select("from_uid as id, count(*) as cnt").Where(where).Group("from_uid").Scan(&subscribings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	subscribers := make([]*_IdCntScanResult, 0)
	where = s.common.BuildOrExpr("to_uid", uids)
	rdb = s.table().Select("to_uid as id, count(*) as cnt").Where(where).Group("to_uid").Scan(&subscribers)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64][2]int32, len(uids))
	for _, r := range subscribings {
		bucket[r.Id] = [2]int32{r.Cnt, 0}
	}
	for _, r := range subscribers {
		if arr, ok := bucket[r.Id]; !ok {
			bucket[r.Id] = [2]int32{0, r.Cnt}
		} else {
			bucket[r.Id] = [2]int32{arr[0], r.Cnt}
		}
	}

	out := make([]*[2]int32, len(uids))
	for idx, uid := range uids {
		arr, ok := bucket[uid]
		if ok {
			out[idx] = &arr
		} else {
			out[idx] = &[2]int32{0, 0}
		}
	}
	return out, nil
}
