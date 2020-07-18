package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/database"
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

func (v *VideoService) WrapVideo(video *po.Video) {
	out := v.db.QueryFirstHelper(&po.User{}, &po.User{Uid: video.AuthorUid})
	if out != nil {
		video.Author = out.(*po.User)
	} else {
		video.Author = nil
	}
}

func (v *VideoService) QueryAll(pageOrder *param.PageOrderParam) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := v.db.QueryMultiHelper(&po.Video{}, pageOrder.Limit, pageOrder.Page, &po.Video{}, v._orderByFunc(pageOrder.Order), &videos)
	for idx := range videos {
		v.WrapVideo(videos[idx])
	}
	return videos, total
}

func (v *VideoService) QueryByUid(uid int32, pageOrder *param.PageOrderParam) ([]*po.Video, int32, database.DbStatus) {
	author := v.userService.QueryByUid(uid)
	if author == nil {
		return nil, 0, database.DbNotFound
	}
	videos := make([]*po.Video, 0)
	total := v.db.QueryMultiHelper(&po.Video{}, pageOrder.Limit, pageOrder.Page, &po.Video{AuthorUid: uid}, v._orderByFunc(pageOrder.Order), &videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, total, database.DbSuccess
}

func (v *VideoService) QueryCountByUid(uid int32) (int32, database.DbStatus) {
	if !v.userService.Exist(uid) {
		return 0, database.DbNotFound
	}
	cnt := v.db.CountHelper(&po.Video{}, &po.Video{AuthorUid: uid})
	return cnt, database.DbSuccess
}

func (v *VideoService) QueryByVid(vid int32) *po.Video {
	out := v.db.QueryFirstHelper(&po.Video{}, &po.Video{Vid: vid})
	if out == nil {
		return nil
	}
	video := out.(*po.Video)
	v.WrapVideo(video)
	return video
}

func (v *VideoService) Exist(vid int32) bool {
	return v.db.ExistHelper(&po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoService) Insert(video *po.Video) database.DbStatus {
	return v.db.InsertHelper(&po.Video{}, video)
}

func (v *VideoService) Update(video *po.Video) database.DbStatus {
	return v.db.UpdateHelper(&po.Video{}, video)
}

func (v *VideoService) Delete(vid int32) database.DbStatus {
	return v.db.DeleteHelper(&po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoService) DeleteBy2Id(vid int32, uid int32) database.DbStatus {
	return v.db.DeleteHelper(&po.Video{}, &po.Video{Vid: vid, AuthorUid: uid})
}
