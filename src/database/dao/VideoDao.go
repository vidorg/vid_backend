package dao

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type VideoDao struct {
	Config  *config.ServerConfig `di:"~"`
	Db      *gorm.DB             `di:"~"`
	UserDao *UserDao             `di:"~"`

	PageSize int32 `di:"-"`
}

func NewVideoDao(dic *xdi.DiContainer) *VideoDao {
	repo := &VideoDao{}
	if !dic.Inject(repo) {
		panic("Inject failed")
	}
	repo.PageSize = repo.Config.MySqlConfig.PageSize
	return repo
}

func (v *VideoDao) WrapVideo(video *po.Video) {
	out := QueryHelper(v.Db, &po.User{}, &po.User{Uid: video.AuthorUid})
	if out != nil {
		video.Author = out.(*po.User)
	} else {
		video.Author = nil
	}
}

func (v *VideoDao) QueryAll(page int32) ([]*po.Video, int32) {
	videos := make([]*po.Video, 0)
	total := PageHelper(v.Db, &po.Video{}, v.PageSize, page, &po.Video{}, &videos)
	for idx := range videos {
		v.WrapVideo(videos[idx])
	}
	return videos, total
}

func (v *VideoDao) QueryByUid(uid int32, page int32) ([]*po.Video, int32, database.DbStatus) {
	author := v.UserDao.QueryByUid(uid)
	if author == nil {
		return nil, 0, database.DbNotFound
	}
	videos := make([]*po.Video, 0)
	total := PageHelper(v.Db, &po.Video{}, v.PageSize, page, &po.Video{AuthorUid: uid}, &videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, total, database.DbSuccess
}

func (v *VideoDao) QueryCountByUid(uid int32) (int32, database.DbStatus) {
	if !v.UserDao.Exist(uid) {
		return 0, database.DbNotFound
	}
	cnt := CountHelper(v.Db, &po.Video{}, &po.Video{AuthorUid: uid})
	return cnt, database.DbSuccess
}

func (v *VideoDao) QueryByVid(vid int32) *po.Video {
	out := QueryHelper(v.Db, &po.Video{}, &po.Video{Vid: vid})
	if out == nil {
		return nil
	}
	video := out.(*po.Video)
	v.WrapVideo(video)
	return video
}

func (v *VideoDao) Exist(vid int32) bool {
	return ExistHelper(v.Db, &po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoDao) Insert(video *po.Video) database.DbStatus {
	return InsertHelper(v.Db, &po.Video{}, video)
}

func (v *VideoDao) Update(video *po.Video) database.DbStatus {
	return UpdateHelper(v.Db, &po.Video{}, video)
}

func (v *VideoDao) Delete(vid int32) database.DbStatus {
	return DeleteHelper(v.Db, &po.Video{}, &po.Video{Vid: vid})
}

func (v *VideoDao) DeleteBy2Id(vid int32, uid int32) database.DbStatus {
	return DeleteHelper(v.Db, &po.Video{}, &po.Video{Vid: vid, AuthorUid: uid})
}
