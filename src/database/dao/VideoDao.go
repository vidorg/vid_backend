package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type VideoDao struct {
	db       *gorm.DB
	pageSize int
	userDao  *UserDao
}

func VideoRepository(config *config.DatabaseConfig) *VideoDao {
	return &VideoDao{
		db:       database.SetupDBConn(config),
		pageSize: config.PageSize,
		userDao:  UserRepository(config),
	}
}

func (v *VideoDao) QueryAll(page int) (videos []*po.Video, count int) {
	v.db.Model(&po.Video{}).Count(&count)
	v.db.Model(&po.Video{}).Limit(v.pageSize).Offset((page - 1) * v.pageSize).Find(&videos)
	for idx := range videos {
		author := &po.User{}
		v.db.Model(&po.User{}).Where(&po.User{Uid: videos[idx].AuthorUid}).Find(author)
		videos[idx].Author = author
	}
	return videos, count
}

func (v *VideoDao) QueryByUid(uid int, page int) (videos []*po.Video, count int) {
	author := v.userDao.QueryByUid(uid)
	if author == nil {
		return nil, 0
	}
	video := &po.Video{AuthorUid: uid}
	v.db.Model(&po.Video{}).Where(video).Count(&count)
	v.db.Model(&po.Video{}).Limit(v.pageSize).Offset((page - 1) * v.pageSize).Where(video).Find(&videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, count
}

func (v *VideoDao) QueryCount(uid int) (int, database.DbStatus) {
	if !v.userDao.Exist(uid) {
		return 0, database.DbNotFound
	}
	video := &po.Video{AuthorUid: uid}
	var count int
	v.db.Model(&po.Video{}).Where(video).Count(&count)
	return count, database.DbSuccess
}

func (v *VideoDao) QueryByVid(vid int) *po.Video {
	video := &po.Video{Vid: vid}
	rdb := v.db.Model(&po.Video{}).Where(video).First(video)
	if rdb.RecordNotFound() {
		return nil
	}
	user := &po.User{Uid: video.AuthorUid}
	rdb = v.db.Model(&po.User{}).Where(user).First(user)
	if rdb.RecordNotFound() {
		return nil
	}
	video.Author = user
	return video
}

func (v *VideoDao) Exist(vid int) bool {
	video := &po.Video{Vid: vid}
	rdb := v.db.Model(&po.Video{}).Where(video)
	return !rdb.RecordNotFound()
}

func (v *VideoDao) Insert(video *po.Video) database.DbStatus {
	rdb := v.db.Model(&po.Video{}).Create(video)
	if rdb.Error != nil {
		if database.IsDuplicateError(rdb.Error) {
			return database.DbExisted
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}

func (v *VideoDao) Update(video *po.Video, uid int) database.DbStatus {
	if video.AuthorUid != uid {
		return database.DbNotFound
	}
	rdb := v.db.Model(&po.Video{}).Update(video)
	if rdb.Error != nil {
		if database.IsNotFoundError(rdb.Error) {
			return database.DbNotFound
		} else if database.IsDuplicateError(rdb.Error) {
			return database.DbExisted
		} else {
			return database.DbFailed
		}
	}
	return database.DbSuccess
}

func (v *VideoDao) Delete(vid int, uid int) database.DbStatus {
	video := &po.Video{Vid: vid, AuthorUid: uid}
	rdb := v.db.Model(video).Delete(video)
	if rdb != nil {
		if database.IsNotFoundError(rdb.Error) {
			return database.DbNotFound
		}
		return database.DbFailed
	}
	return database.DbSuccess
}
