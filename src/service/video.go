package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type VideoService struct {
	db          *gorm.DB
	userService *UserService

	_orderByFunc func(string) string
}

func NewVideoService() *VideoService {
	return &VideoService{
		db:           xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		userService:  xdi.GetByNameForce(sn.SUserService).(*UserService),
		_orderByFunc: xproperty.GetMapperDefault(&dto.VideoDto{}, &po.Video{}).ApplyOrderBy,
	}
}

func (v *VideoService) QueryAll(pp *param.PageOrderParam) (videos []*po.Video, total int32) {
	total = 0
	v.db.Model(&po.Video{}).Count(&total)

	videos = make([]*po.Video, 0)
	xgorm.WithDB(v.db).Pagination(pp.Limit, pp.Page).
		Model(&po.Video{}).
		Order(v._orderByFunc(pp.Order)).
		Find(&videos)

	for _, video := range videos {
		user := &po.User{}
		rdb := v.db.Model(&po.User{}).Where(&po.User{Uid: video.AuthorUid}).First(user)
		if !rdb.RecordNotFound() {
			video.Author = user
		}
	}

	return videos, total
}

func (v *VideoService) QueryByUid(uid int32, pp *param.PageOrderParam) (videos []*po.Video, total int32, status xstatus.DbStatus) {
	author := v.userService.QueryByUid(uid)
	if author == nil {
		return nil, 0, xstatus.DbNotFound
	}

	total = 0
	v.db.Model(&po.Video{}).Where(&po.Video{AuthorUid: uid}).Count(&total)

	videos = make([]*po.Video, 0)
	xgorm.WithDB(v.db).Pagination(pp.Limit, pp.Page).
		Model(&po.Video{}).
		Order(v._orderByFunc(pp.Order)).
		Where(&po.Video{AuthorUid: uid}).
		Find(&videos)

	for idx := range videos {
		videos[idx].Author = author
	}

	return videos, total, xstatus.DbSuccess
}

func (v *VideoService) QueryCountByUid(uid int32) (int32, xstatus.DbStatus) {
	if !v.userService.Exist(uid) {
		return 0, xstatus.DbNotFound
	}

	total := 0
	v.db.Model(&po.Video{}).Where(&po.Video{AuthorUid: uid}).Count(&total)

	return int32(total), xstatus.DbSuccess
}

func (v *VideoService) QueryByVid(vid int32) *po.Video {
	video := &po.Video{}
	rdb := v.db.Model(&po.Video{}).Where(&po.Video{Vid: vid}).First(video)
	if rdb.RecordNotFound() {
		return nil
	}

	user := &po.User{}
	rdb = v.db.Model(&po.User{}).Where(&po.User{Uid: video.AuthorUid}).First(&user)
	if !rdb.RecordNotFound() {
		video.Author = user
	}

	return video
}

func (v *VideoService) Exist(vid int32) bool {
	return xgorm.WithDB(v.db).Exist(&po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoService) Insert(video *po.Video) xstatus.DbStatus {
	return xgorm.WithDB(v.db).Insert(&po.Video{}, video)
}

func (v *VideoService) Update(video *po.Video) xstatus.DbStatus {
	return xgorm.WithDB(v.db).Update(&po.Video{}, nil, video)
}

func (v *VideoService) Delete(vid int32) xstatus.DbStatus {
	return xgorm.WithDB(v.db).Delete(&po.Video{}, nil, &po.Video{Vid: vid})
}

func (v *VideoService) DeleteBy2Id(vid int32, uid int32) xstatus.DbStatus {
	return xgorm.WithDB(v.db).Delete(&po.Video{}, nil, &po.Video{Vid: vid, AuthorUid: uid})
}
