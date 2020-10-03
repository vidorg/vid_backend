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

type FavoriteService struct {
	db             *gorm.DB
	common         *CommonService
	userService    *UserService
	videoService   *VideoService
	orderbyService *OrderbyService
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		db:             xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		common:         xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		userService:    xdi.GetByNameForce(sn.SUserService).(*UserService),
		videoService:   xdi.GetByNameForce(sn.SVideoService).(*VideoService),
		orderbyService: xdi.GetByNameForce(sn.SOrderbyService).(*OrderbyService),
	}
}

func (f *FavoriteService) table() *gorm.DB {
	return f.db.Table("tbl_favorite")
}

func (f *FavoriteService) favoriteAsso(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Favorites") // association pattern
}

func (f *FavoriteService) favoredAsso(db *gorm.DB, vid uint64) *gorm.Association {
	return db.Model(&po.Video{Vid: vid}).Association("Favoreds") // association pattern
}

func (f *FavoriteService) QueryFavorites(uid uint64, pp *param.PageOrderParam) ([]*po.Video, int32, error) {
	ok, err := f.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(f.favoriteAsso(f.db, uid).Count())
	videos := make([]*po.Video, 0)
	err = f.favoriteAsso(xgorm.WithDB(f.db).Pagination(pp.Limit, pp.Page).Order(f.orderbyService.FavoriteForVideo(pp.Order)), uid).Find(&videos)
	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

func (f *FavoriteService) QueryFavoreds(vid uint64, pp *param.PageOrderParam) ([]*po.User, int32, error) {
	ok, err := f.videoService.Existed(vid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(f.favoredAsso(f.db, vid).Count())
	users := make([]*po.User, 0)
	err = f.favoredAsso(xgorm.WithDB(f.db).Pagination(pp.Limit, pp.Page).Order(f.orderbyService.FavoriteForUser(pp.Order)), vid).Find(&users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (f *FavoriteService) InsertFavorite(uid uint64, vid uint64) (xstatus.DbStatus, error) {
	ok1, err1 := f.userService.Existed(uid)
	ok2, err2 := f.videoService.Existed(vid)
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
	rdb := f.table().Where("uid = ? AND vid = ?", uid, vid).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt > 0 {
		return xstatus.DbExisted, nil // existed
	}

	err := f.favoriteAsso(f.db, uid).Append(&po.Video{Vid: vid})
	if err != nil {
		return xstatus.DbFailed, err
	}

	return xstatus.DbSuccess, nil
}

func (f *FavoriteService) DeleteFavorite(uid uint64, vid uint64) (xstatus.DbStatus, error) {
	ok1, err1 := f.userService.Existed(uid)
	ok2, err2 := f.videoService.Existed(vid)
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
	rdb := f.table().Where("uid = ? AND vid = ?", uid, vid).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt == 0 {
		return xstatus.DbTagA, nil // not found
	}

	err := f.favoriteAsso(f.db, uid).Delete(&po.Video{Vid: vid})
	if err != nil {
		return xstatus.DbFailed, err
	}

	return xstatus.DbSuccess, nil
}

func (f *FavoriteService) CheckFavorite(uid uint64, vids []uint64) ([]bool, error) {
	if len(vids) == 0 {
		return []bool{}, nil
	}

	favorites := make([]*_IdScanResult, 0)
	where := f.common.BuildOrExpr("vid", vids)
	rdb := f.table().Select("vid as id").Where("uid = ?", uid).Where(where).Scan(&favorites)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]bool, len(vids))
	for _, r := range favorites {
		bucket[r.Id] = true
	}

	out := make([]bool, len(vids))
	for idx, vid := range vids {
		_, ok := bucket[vid]
		out[idx] = ok
	}
	return out, nil
}

func (f *FavoriteService) QueryFavoriteCount(uids []uint64) ([]int32, error) {
	if len(uids) == 0 {
		return []int32{}, nil
	}

	counts := make([]*_IdCntScanResult, 0)
	where := f.common.BuildOrExpr("uid", uids)
	rdb := f.table().Select("uid as id, count(*) as cnt").Where(where).Group("uid").Scan(&counts)
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

func (f *FavoriteService) QueryFavoredCount(vids []uint64) ([]int32, error) {
	if len(vids) == 0 {
		return []int32{}, nil
	}

	counts := make([]*_IdCntScanResult, 0)
	where := f.common.BuildOrExpr("vid", vids)
	rdb := f.table().Select("vid as id, count(*) as cnt").Where(where).Group("vid").Scan(&counts)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]int32)
	for _, cnt := range counts {
		bucket[cnt.Id] = cnt.Cnt
	}

	out := make([]int32, len(vids))
	for idx, vid := range vids {
		if cnt, ok := bucket[vid]; ok {
			out[idx] = cnt
		}
	}
	return out, nil
}
