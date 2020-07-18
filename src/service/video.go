package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
)

type VideoService struct {
	Db          *database.GormHelper `di:"~"`
	Logger      *logrus.Logger       `di:"~"`
	UserService *UserService         `di:"~"`

	OrderByFunc func(string) string `di:"-"`
}

func NewVideoService(dic *xdi.DiContainer) *VideoService {
	repo := &VideoService{}
	dic.MustInject(repo)
	repo.OrderByFunc = xproperty.GetMapperDefault(&dto.VideoDto{}, &po.Video{}).ApplyOrderBy
	return repo
}

func (v *VideoService) WrapVideo(video *po.Video) {
	out := v.Db.QueryFirstHelper(&po.User{}, &po.User{Uid: video.AuthorUid})
	if out != nil {
		video.Author = out.(*po.User)
	} else {
		video.Author = nil
	}
}

func (v *VideoService) QueryAll(pageOrder *param.PageOrderParam) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := v.Db.QueryMultiHelper(&po.Video{}, pageOrder.Limit, pageOrder.Page, &po.Video{}, v.OrderByFunc(pageOrder.Order), &videos)
	for idx := range videos {
		v.WrapVideo(videos[idx])
	}
	return videos, total
}

func (v *VideoService) QueryByUid(uid int32, pageOrder *param.PageOrderParam) ([]*po.Video, int32, database.DbStatus) {
	author := v.UserService.QueryByUid(uid)
	if author == nil {
		return nil, 0, database.DbNotFound
	}
	videos := make([]*po.Video, 0)
	total := v.Db.QueryMultiHelper(&po.Video{}, pageOrder.Limit, pageOrder.Page, &po.Video{AuthorUid: uid}, v.OrderByFunc(pageOrder.Order), &videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, total, database.DbSuccess
}

func (v *VideoService) QueryCountByUid(uid int32) (int32, database.DbStatus) {
	if !v.UserService.Exist(uid) {
		return 0, database.DbNotFound
	}
	cnt := v.Db.CountHelper(&po.Video{}, &po.Video{AuthorUid: uid})
	return cnt, database.DbSuccess
}

func (v *VideoService) QueryByVid(vid int32) *po.Video {
	out := v.Db.QueryFirstHelper(&po.Video{}, &po.Video{Vid: vid})
	if out == nil {
		return nil
	}
	video := out.(*po.Video)
	v.WrapVideo(video)
	return video
}

func (v *VideoService) Exist(vid int32) bool {
	return v.Db.ExistHelper(&po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoService) Insert(video *po.Video) database.DbStatus {
	return v.Db.InsertHelper(&po.Video{}, video)
}

func (v *VideoService) Update(video *po.Video) database.DbStatus {
	return v.Db.UpdateHelper(&po.Video{}, video)
}

func (v *VideoService) Delete(vid int32) database.DbStatus {
	return v.Db.DeleteHelper(&po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoService) DeleteBy2Id(vid int32, uid int32) database.DbStatus {
	return v.Db.DeleteHelper(&po.Video{}, &po.Video{Vid: vid, AuthorUid: uid})
}
