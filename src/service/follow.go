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

type FollowService struct {
	db             *gorm.DB
	common         *CommonService
	userService    *UserService
	orderbyService *OrderbyService
}

func NewFollowService() *FollowService {
	return &FollowService{
		db:             xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		common:         xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		userService:    xdi.GetByNameForce(sn.SUserService).(*UserService),
		orderbyService: xdi.GetByNameForce(sn.SOrderbyService).(*OrderbyService),
	}
}

func (s *FollowService) table() *gorm.DB {
	return s.db.Table("tbl_follow")
}

func (s *FollowService) followerAsso(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Followers") // association pattern
}

func (s *FollowService) followingAsso(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Followings") // association pattern
}

func (s *FollowService) QueryFollowers(uid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.followerAsso(s.db, uid).Count())
	users := make([]*po.User, 0)
	err = s.followerAsso(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderbyService.FollowForUser(pp.Order)), uid).Find(&users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *FollowService) QueryFollowings(uid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := s.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(s.followingAsso(s.db, uid).Count())
	users := make([]*po.User, 0)
	err = s.followingAsso(xgorm.WithDB(s.db).Pagination(pp.Limit, pp.Page).Order(s.orderbyService.FollowForUser(pp.Order)), uid).Find(&users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *FollowService) FollowUser(uid uint64, to uint64) (xstatus.DbStatus, error) {
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

	err := s.followerAsso(s.db, to).Append(&po.User{Uid: uid})
	if err != nil {
		return xstatus.DbFailed, err
	}

	return xstatus.DbSuccess, nil
}

func (s *FollowService) UnfollowUser(uid uint64, to uint64) (xstatus.DbStatus, error) {
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

	err := s.followerAsso(s.db, to).Delete(&po.User{Uid: uid})
	if err != nil {
		return xstatus.DbFailed, err
	}

	return xstatus.DbSuccess, nil
}

func (s *FollowService) CheckFollow(me uint64, uids []uint64) (ingerPairs []*[2]bool, err error) {
	if len(uids) == 0 {
		return []*[2]bool{}, nil
	}

	followings := make([]*_IdScanResult, 0)
	where := s.common.BuildOrExpr("to_uid", uids)
	// TODO tx.Statement.Selects = []string{strings.Join(fields, " ")}
	rdb := s.table().Select("to_uid as id").Where("from_uid = ?", me).Where(where).Scan(&followings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	followers := make([]*_IdScanResult, 0)
	where = s.common.BuildOrExpr("from_uid", uids)
	rdb = s.table().Select("from_uid as id").Where("to_uid = ?", me).Where(where).Scan(&followers)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64][2]bool, len(uids))
	for _, r := range followings {
		bucket[r.Id] = [2]bool{true, false}
	}
	for _, r := range followers {
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

func (s *FollowService) QueryFollowCount(uids []uint64) (ingerPairs []*[2]int32, err error) {
	if len(uids) == 0 {
		return []*[2]int32{}, nil
	}

	followings := make([]*_IdCntScanResult, 0)
	where := s.common.BuildOrExpr("from_uid", uids)
	rdb := s.table().Select("from_uid as id, count(*) as cnt").Where(where).Group("from_uid").Scan(&followings)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	followers := make([]*_IdCntScanResult, 0)
	where = s.common.BuildOrExpr("to_uid", uids)
	rdb = s.table().Select("to_uid as id, count(*) as cnt").Where(where).Group("to_uid").Scan(&followers)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64][2]int32, len(uids))
	for _, r := range followings {
		bucket[r.Id] = [2]int32{r.Cnt, 0}
	}
	for _, r := range followers {
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
