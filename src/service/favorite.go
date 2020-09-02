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
	tableName    string
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		db:           xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService:  xdi.GetByNameForce(sn.SUserService).(*UserService),
		videoService: xdi.GetByNameForce(sn.SVideoService).(*VideoService),
		common:       xdi.GetByNameForce(sn.SCommonService).(*CommonService),
		userOrderBy:  xgorm.OrderByFunc(dto.BuildUserPropertyMapper()),
		videoOrderBy: xgorm.OrderByFunc(dto.BuildVideoPropertyMapper()),
		tableName:    "tbl_favorite",
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
	// TODO
	return xstatus.DbSuccess, nil
}

func (f *FavoriteService) DeleteFavorite(uid uint64, vid uint64) (xstatus.DbStatus, error) {
	// TODO
	return xstatus.DbSuccess, nil
}

func (f *FavoriteService) QueryCountByVids(vids []uint64) ([]int32, error) {
	// TODO
	return nil, nil
}

func (f *FavoriteService) QueryCountByUids(uids []uint64) ([]int32, error) {
	// TODO
	return nil, nil
}

func (f *FavoriteService) CheckFavorites(uid uint64, vids []uint64) ([]bool, error) {
	// TODO
	return nil, nil
}
