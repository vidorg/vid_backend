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

type FavoriteService struct {
	db           *gorm.DB
	userService  *UserService
	videoService *VideoService
	common       *CommonService
	userOrderBy  func(string) string
	videoOrderBy func(string) string
	tblName      string
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		db:           xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService:  xdi.GetByNameForce(sn.SUserService).(*UserService),
		videoService: xdi.GetByNameForce(sn.SVideoService).(*VideoService),
		common:       xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		userOrderBy:  xgorm.OrderByFunc(dto.BuildUserPropertyMapper()),
		videoOrderBy: xgorm.OrderByFunc(dto.BuildVideoPropertyMapper()),
		tblName:      "tbl_favorite",
	}
}

func (f *FavoriteService) favoriteAssoDB(db *gorm.DB, uid uint64) *gorm.Association {
	return db.Model(&po.User{Uid: uid}).Association("Favorites") // association pattern
}

func (f *FavoriteService) favoriteAsso(uid uint64) *gorm.Association {
	return f.favoriteAssoDB(f.db, uid)
}

func (f *FavoriteService) favoredAssoDB(db *gorm.DB, vid uint64) *gorm.Association {
	return db.Model(&po.Video{Vid: vid}).Association("Favoreds") // association pattern
}

func (f *FavoriteService) favoredAsso(uid uint64) *gorm.Association {
	return f.favoredAssoDB(f.db, uid)
}

func (f *FavoriteService) QueryFavorites(uid uint64, pp *param.PageOrderParam) ([]*po.Video, int32, error) {
	ok, err := f.userService.Existed(uid)
	if err != nil {
		return nil, 0, err
	} else if !ok {
		return nil, 0, nil
	}

	total := int32(f.favoriteAsso(uid).Count())
	videos := make([]*po.Video, 0)
	rac := f.favoriteAssoDB(xgorm.WithDB(f.db).Pagination(pp.Limit, pp.Page).Order(f.videoOrderBy(pp.Order)), uid).Find(&videos)
	if rac.Error != nil {
		return nil, 0, rac.Error
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

	total := int32(f.favoredAsso(vid).Count())
	users := make([]*po.User, 0)
	rac := f.favoredAssoDB(xgorm.WithDB(f.db).Pagination(pp.Limit, pp.Page).Order(f.userOrderBy(pp.Order)), vid).Find(&users)
	if rac.Error != nil {
		return nil, 0, rac.Error
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

	cnt := 0
	rdb := f.db.Table(f.tblName).Where("uid = ? AND vid = ?", uid, vid).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt > 0 {
		return xstatus.DbExisted, nil
	}

	ras := f.favoriteAsso(uid).Append(&po.Video{Vid: vid})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
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

	cnt := 0
	rdb := f.db.Table(f.tblName).Where("uid = ? AND vid = ?", uid, vid).Count(&cnt)
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if cnt == 0 {
		return xstatus.DbTagA, nil
	}

	ras := f.favoriteAsso(uid).Delete(&po.Video{Vid: vid})
	if ras.Error != nil {
		return xstatus.DbFailed, ras.Error
	}

	return xstatus.DbSuccess, nil
}

func (f *FavoriteService) QueryFavoredCount(uids []uint64) ([]int32, error) {
	if len(uids) == 0 {
		return []int32{}, nil
	}

	counts := make([]*_IdCntPair, 0)
	where := f.common.BuildOrExp("uid", uids)
	rdb := f.db.Table(f.tblName).Select("uid as id, count(*) as cnt").Where(where).Group("uid").Scan(&counts)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]int32)
	for _, cnt := range counts {
		bucket[cnt.Id] = cnt.Cnt
	}
	out := make([]int32, len(uids))
	for idx, uid := range uids {
		cnt, ok := bucket[uid]
		if ok {
			out[idx] = cnt
		}
	}
	return out, nil
}

func (f *FavoriteService) QueryFavoriteCount(vids []uint64) ([]int32, error) {
	if len(vids) == 0 {
		return []int32{}, nil
	}

	counts := make([]*_IdCntPair, 0)
	where := f.common.BuildOrExp("vid", vids)
	rdb := f.db.Table(f.tblName).Select("vid as id, count(*) as cnt").Where(where).Group("vid").Scan(&counts)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]int32)
	for _, cnt := range counts {
		bucket[cnt.Id] = cnt.Cnt
	}
	out := make([]int32, len(vids))
	for idx, vid := range vids {
		cnt, ok := bucket[vid]
		if ok {
			out[idx] = cnt
		}
	}
	return out, nil
}

func (f *FavoriteService) CheckFavorite(uid uint64, vids []uint64) ([]bool, error) {
	if len(vids) == 0 {
		return []bool{}, nil
	}

	favorites := make([]*_UidVidPair, 0)
	where := f.common.BuildOrExp("vid", vids)
	rdb := f.db.Table(f.tblName).Select("uid, vid").Where("uid = ?", uid).Where(where).Group("uid, vid").Scan(&favorites)
	if rdb.Error != nil {
		return nil, rdb.Error
	}

	bucket := make(map[uint64]bool, len(vids))
	for _, pair := range favorites {
		bucket[pair.Vid] = true
	}
	out := make([]bool, len(vids))
	for idx, vid := range vids {
		video, ok := bucket[vid]
		if ok {
			out[idx] = video
		}
	}
	return out, nil
}
