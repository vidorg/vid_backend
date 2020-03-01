package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
)

type VideoDao struct {
	Config          *config.ServerConfig       `di:"~"`
	Db              *database.DbHelper         `di:"~"`
	PropertyMappers *xproperty.PropertyMappers `di:"~"`
	UserDao         *UserDao                   `di:"~"`

	PageSize    int32               `di:"-"`
	OrderByFunc func(string) string `di:"-"`
}

func NewVideoDao(dic *xdi.DiContainer) *VideoDao {
	repo := &VideoDao{}
	if !dic.Inject(repo) {
		log.Fatalln("Inject failed")
	}
	repo.PageSize = repo.Config.MySqlConfig.PageSize
	repo.OrderByFunc = repo.PropertyMappers.GetPropertyMapping(&dto.VideoDto{}, &po.Video{}).ApplyOrderBy
	return repo
}

func (v *VideoDao) WrapVideo(video *po.Video) {
	out := v.Db.QueryHelper(&po.User{}, &po.User{Uid: video.AuthorUid})
	if out != nil {
		video.Author = out.(*po.User)
	} else {
		video.Author = nil
	}
}

func (v *VideoDao) QueryAll(page int32, orderBy string) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := v.Db.QueryMultiHelper(&po.Video{}, v.PageSize, page, &po.Video{}, v.OrderByFunc(orderBy), &videos)
	for idx := range videos {
		v.WrapVideo(videos[idx])
	}
	return videos, total
}

func (v *VideoDao) QueryByUid(uid int32, page int32, orderBy string) ([]*po.Video, int32, database.DbStatus) {
	author := v.UserDao.QueryByUid(uid)
	if author == nil {
		return nil, 0, database.DbNotFound
	}
	videos := make([]*po.Video, 0)
	total := v.Db.QueryMultiHelper(&po.Video{}, v.PageSize, page, &po.Video{AuthorUid: uid}, v.OrderByFunc(orderBy), &videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, total, database.DbSuccess
}

func (v *VideoDao) QueryCountByUid(uid int32) (int32, database.DbStatus) {
	if !v.UserDao.Exist(uid) {
		return 0, database.DbNotFound
	}
	cnt := v.Db.CountHelper(&po.Video{}, &po.Video{AuthorUid: uid})
	return cnt, database.DbSuccess
}

func (v *VideoDao) QueryByVid(vid int32) *po.Video {
	out := v.Db.QueryHelper(&po.Video{}, &po.Video{Vid: vid})
	if out == nil {
		return nil
	}
	video := out.(*po.Video)
	v.WrapVideo(video)
	return video
}

func (v *VideoDao) Exist(vid int32) bool {
	return v.Db.ExistHelper(&po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoDao) Insert(video *po.Video) database.DbStatus {
	return v.Db.InsertHelper(&po.Video{}, video)
}

func (v *VideoDao) Update(video *po.Video) database.DbStatus {
	return v.Db.UpdateHelper(&po.Video{}, video)
}

func (v *VideoDao) Delete(vid int32) database.DbStatus {
	return v.Db.DeleteHelper(&po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoDao) DeleteBy2Id(vid int32, uid int32) database.DbStatus {
	return v.Db.DeleteHelper(&po.Video{}, &po.Video{Vid: vid, AuthorUid: uid})
}
