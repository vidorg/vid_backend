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

func NewVideoDao(dic xdi.DiContainer) *VideoDao {
	repo := &VideoDao{}
	dic.Inject(repo)
	if xdi.HasNilDi(repo) {
		panic("Has nil di field")
	}

	repo.PageSize = repo.Config.MySqlConfig.PageSize
	return repo
}

func (v *VideoDao) QueryAll(page int32) (videos []*po.Video, count int32) {
	v.Db.Model(&po.Video{}).Count(&count)
	v.Db.Model(&po.Video{}).Limit(v.PageSize).Offset((page - 1) * v.PageSize).Find(&videos)
	for idx := range videos {
		author := &po.User{}
		v.Db.Where(&po.User{Uid: videos[idx].AuthorUid}).Find(author)
		videos[idx].Author = author
	}
	return videos, count
}

func (v *VideoDao) QueryByUid(uid int32, page int32) (videos []*po.Video, count int32, status database.DbStatus) {
	author := v.UserDao.QueryByUid(uid)
	if author == nil {
		return nil, 0, database.DbNotFound
	}
	video := &po.Video{AuthorUid: uid}
	v.Db.Model(&po.Video{}).Where(video).Count(&count)
	v.Db.Model(&po.Video{}).Limit(v.PageSize).Offset((page - 1) * v.PageSize).Where(video).Find(&videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, count, database.DbSuccess
}

func (v *VideoDao) QueryCount(uid int32) (count int32, status database.DbStatus) {
	if !v.UserDao.Exist(uid) {
		return 0, database.DbNotFound
	}
	video := &po.Video{AuthorUid: uid}
	v.Db.Model(&po.Video{}).Where(video).Count(&count)
	return count, database.DbSuccess
}

func (v *VideoDao) QueryByVid(vid int32) *po.Video {
	video := &po.Video{Vid: vid}
	rdb := v.Db.Model(&po.Video{}).Where(video).First(video)
	if rdb.RecordNotFound() {
		return nil
	}
	user := &po.User{Uid: video.AuthorUid}
	rdb = v.Db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		// nullable
		user = &po.User{Uid: -1}
	}
	video.Author = user
	return video
}

func (v *VideoDao) Exist(vid int32) bool {
	video := &po.Video{Vid: vid}
	cnt := 0
	v.Db.Model(&po.Video{}).Where(video).Count(&cnt)
	return cnt > 0
}

func (v *VideoDao) Insert(video *po.Video) database.DbStatus {
	rdb := v.Db.Model(&po.Video{}).Model(&po.Video{}).Create(video)
	if database.IsDuplicateError(rdb.Error) {
		return database.DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (v *VideoDao) Update(video *po.Video) database.DbStatus {
	rdb := v.Db.Model(&po.Video{}).Update(video)
	if rdb.Error != nil {
		if database.IsDuplicateError(rdb.Error) {
			return database.DbExisted
		} else {
			return database.DbFailed
		}
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}

func (v *VideoDao) Delete(vid int32) database.DbStatus {
	video := &po.Video{Vid: vid}
	rdb := v.Db.Model(&po.Video{}).Delete(video)
	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}
