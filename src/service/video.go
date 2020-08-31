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

type VideoService struct {
	db          *gorm.DB
	userService *UserService
	orderBy     func(string) string
}

func NewVideoService() *VideoService {
	return &VideoService{
		db:          xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService: xdi.GetByNameForce(sn.SUserService).(*UserService),
		orderBy:     xgorm.OrderByFunc(xproperty.GetDefaultMapper(&dto.VideoDto{}, &po.Video{}).GetDict()),
	}
}

func (v *VideoService) QueryAll(pp *param.PageOrderParam) ([]*po.Video, int32, error) {
	total := int32(0)
	v.db.Model(&po.Video{}).Count(&total)

	videos := make([]*po.Video, 0)
	rdb := xgorm.WithDB(v.db).Pagination(pp.Limit, pp.Page).Model(&po.User{}).Order(v.orderBy(pp.Order)).Find(&videos)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	// TODO Use OR operator
	for _, video := range videos {
		user := &po.User{}
		rdb := v.db.Model(&po.User{}).Where(&po.User{Uid: video.AuthorUid}).First(user) // TODO Use Related
		if rdb.RecordNotFound() {
			video.Author = nil
		} else if rdb.Error != nil {
			return nil, 0, rdb.Error
		} else {
			video.Author = user
		}
	}

	return videos, total, nil
}

func (v *VideoService) QueryByUid(uid uint64, pp *param.PageOrderParam) ([]*po.Video, int32, error) {
	author, err := v.userService.QueryByUid(uid)
	if err != nil {
		return nil, 0, err
	} else if author == nil {
		return nil, 0, nil
	}

	total := int32(0)
	v.db.Model(&po.Video{}).Where(&po.Video{AuthorUid: uid}).Count(&total)

	videos := make([]*po.Video, 0)
	rdb := xgorm.WithDB(v.db).Pagination(pp.Limit, pp.Page).Model(&po.Video{}).Order(v.orderBy(pp.Order)).Where(&po.Video{AuthorUid: uid}).Find(&videos)
	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}

	for _, video := range videos {
		video.Author = author
	}

	return videos, total, nil
}

func (v *VideoService) QueryByVid(vid uint64) (*po.Video, error) {
	video := &po.Video{}
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).First(video)
	if rdb.RecordNotFound() {
		return nil, nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	}

	user := &po.User{}
	rdb = v.db.Model(&po.User{}).Where(&po.User{Uid: video.AuthorUid}).First(&user)
	if rdb.RecordNotFound() {
		video.Author = nil
	} else if rdb.Error != nil {
		return nil, rdb.Error
	} else {
		video.Author = user
	}

	return video, nil
}

func (v *VideoService) QueryCountByUid(uid uint64) (int32, error) {
	ok, err := v.userService.Existed(uid)
	if err != nil {
		return 0, err
	} else if !ok {
		return -1, nil
	}

	total := int32(0)
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{AuthorUid: uid}).Count(&total)
	if rdb.Error != nil {
		return 0, rdb.Error
	}

	return total, nil
}

func (v *VideoService) Existed(vid uint64) (bool, error) {
	cnt := 0
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).Count(&cnt)
	if rdb.Error != nil {
		return false, rdb.Error
	}

	return cnt > 0, nil
}

func (v *VideoService) Insert(pa *param.InsertVideoParam, uid uint64) (xstatus.DbStatus, error) {
	video := pa.ToPo()
	video.AuthorUid = uid

	rdb := v.db.Model(&po.Video{}).Create(video)
	return xgorm.CreateErr(rdb)
}

func (v *VideoService) Update(vid uint64, video *param.UpdateVideoParam) (xstatus.DbStatus, error) {
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).Update(video.ToMap())
	return xgorm.UpdateErr(rdb)
}

func (v *VideoService) Delete(vid uint64) (xstatus.DbStatus, error) {
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).Delete(&po.Video{Vid: vid})
	return xgorm.DeleteErr(rdb)
}
