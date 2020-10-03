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
	channelService *ChannelService
	orderbyService *OrderbyService
}

func NewSubscribeService() *SubscribeService {
	return &SubscribeService{
		db:             xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		common:         xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		userService:    xdi.GetByNameForce(sn.SUserService).(*UserService),
		channelService: xdi.GetByNameForce(sn.SChannelService).(*ChannelService),
		orderbyService: xdi.GetByNameForce(sn.SOrderbyService).(*OrderbyService),
	}
}

func (s *SubscribeService) table() *gorm.DB {
	return s.db.Table("tbl_subscribe")
}

func (s *SubscribeService) subscribingAsso(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Subscribings") // association pattern
}

func (s *SubscribeService) subscriberAsso(db *gorm.DB, cid uint64) *gorm.Association {
	return db.Model(&po.Channel{Cid: cid}).Association("Subscribers") // association pattern
}

func (s *SubscribeService) QuerySubscribings(uid uint64, pp *param.PageOrderParam) ([]*po.Channel, int32, error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.subscribingAsso(s.db, uid).Count())
	channels := make([]*po.Channel, 0)
	err = s.subscribingAsso(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderbyService.SubscribeForChannel(pp.Order)), uid).Find(&channels)
	if err != nil {
		return nil, 0, err
	}

	return channels, total, nil
}

func (s *SubscribeService) QuerySubscribers(cid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := s.channelService.Existed(cid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.subscriberAsso(s.db, cid).Count())
	users := make([]*po.User, 0)
	err = s.subscriberAsso(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderbyService.SubscribeForUser(pp.Order)), cid).Find(&users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *SubscribeService) SubscribeChannel(uid uint64, cid uint64) (xstatus.DbStatus, error) {
	ok1, err1 := s.userService.Existed(uid)
	ok2, err2 := s.channelService.Existed(cid)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 {
		return xstatus.DbTagB, nil
	} else if !ok2 {
		return xstatus.DbTagC, nil
	}

	cnt := int64(0)
	rdb := s.table().Where("uid = ? AND cid = ?", uid, cid).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt > 0 {
		return xstatus.DbExisted, nil // existed
	}

	err := s.subscribingAsso(s.db, uid).Append(&po.Channel{Cid: cid})
	if err != nil {
		return xstatus.DbFailed, err
	}

	return xstatus.DbSuccess, nil
}

func (s *SubscribeService) UnsubscribeChannel(uid uint64, cid uint64) (xstatus.DbStatus, error) {
	ok1, err1 := s.userService.Existed(uid)
	ok2, err2 := s.channelService.Existed(cid)
	if err1 != nil {
		return xstatus.DbFailed, err1
	} else if err2 != nil {
		return xstatus.DbFailed, err2
	} else if !ok1 {
		return xstatus.DbTagB, nil
	} else if !ok2 {
		return xstatus.DbTagC, nil
	}

	cnt := int64(0)
	rdb := s.table().Where("uid = ? AND cid = ?", uid, cid).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt == 0 {
		return xstatus.DbTagA, nil // not found
	}

	err := s.subscribingAsso(s.db, uid).Delete(&po.Channel{Cid: cid})
	if err != nil {
		return xstatus.DbFailed, err
	}

	return xstatus.DbSuccess, nil
}

func (s *SubscribeService) CheckSubscribe(uid uint64, cids []uint64) ([]bool, error) {
	if len(cids) == 0 {
		return []bool{}, nil
	}

	subscribings := make([]*_IdScanResult, 0)
	where := s.common.BuildOrExpr("cid", cids)
	rdb := s.table().Select("cid as id").Where("uid = ?", uid).Where(where).Scan(&subscribings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]bool, len(cids))
	for _, r := range subscribings {
		bucket[r.Id] = true
	}
	out := make([]bool, len(cids))
	for idx, cid := range cids {
		_, ok := bucket[cid]
		out[idx] = ok
	}
	return out, nil
}

func (s *SubscribeService) QuerySubscribingCount(uids []uint64) ([]int32, error) {
	if len(uids) == 0 {
		return []int32{}, nil
	}

	counts := make([]*_IdCntScanResult, 0)
	where := s.common.BuildOrExpr("uid", uids)
	rdb := s.table().Select("uid as id, count(*) as cnt").Where(where).Group("uid").Scan(&counts)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]int32)
	for _, r := range counts {
		bucket[r.Id] = r.Cnt
	}
	out := make([]int32, len(uids))
	for idx, uid := range uids {
		if cnt, ok := bucket[uid]; ok {
			out[idx] = cnt
		}
	}
	return out, nil
}

func (s *SubscribeService) QuerySubscriberCount(cids []uint64) ([]int32, error) {
	if len(cids) == 0 {
		return []int32{}, nil
	}

	counts := make([]*_IdCntScanResult, 0)
	where := s.common.BuildOrExpr("cid", cids)
	rdb := s.table().Select("cid as id, count(*) as cnt").Where(where).Group("cid").Scan(&counts)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]int32)
	for _, r := range counts {
		bucket[r.Id] = r.Cnt
	}
	out := make([]int32, len(cids))
	for idx, cid := range cids {
		if cnt, ok := bucket[cid]; ok {
			out[idx] = cnt
		}
	}
	return out, nil
}
