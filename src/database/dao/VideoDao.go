package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type VideoDao struct {
	db       *gorm.DB
	pageSize int32
	userDao  *UserDao
}

func VideoRepository(config *config.MySqlConfig) *VideoDao {
	return &VideoDao{
		db:       database.SetupDBConn(config),
		pageSize: config.PageSize,
		userDao:  UserRepository(config),
	}
}

func (v *VideoDao) QueryAll(page int32) (videos []*po.Video, count int32) {
	v.db.Model(&po.Video{}).Count(&count)
	v.db.Model(&po.Video{}).Limit(v.pageSize).Offset((page - 1) * v.pageSize).Find(&videos)
	for idx := range videos {
		author := &po.User{}
		v.db.Where(&po.User{Uid: videos[idx].AuthorUid}).Find(author)
		videos[idx].Author = author
	}
	return videos, count
}

func (v *VideoDao) QueryByUid(uid int32, page int32) (videos []*po.Video, count int32, status database.DbStatus) {
	author := v.userDao.QueryByUid(uid)
	if author == nil {
		return nil, 0, database.DbNotFound
	}
	video := &po.Video{AuthorUid: uid}
	v.db.Model(&po.Video{}).Where(video).Count(&count)
	v.db.Model(&po.Video{}).Limit(v.pageSize).Offset((page - 1) * v.pageSize).Where(video).Find(&videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, count, database.DbSuccess
}

func (v *VideoDao) QueryCount(uid int32) (count int32, status database.DbStatus) {
	if !v.userDao.Exist(uid) {
		return 0, database.DbNotFound
	}
	video := &po.Video{AuthorUid: uid}
	v.db.Model(&po.Video{}).Where(video).Count(&count)
	return count, database.DbSuccess
}

func (v *VideoDao) QueryByVid(vid int32) *po.Video {
	video := &po.Video{Vid: vid}
	rdb := v.db.Model(&po.Video{}).Where(video).First(video)
	if rdb.RecordNotFound() {
		return nil
	}
	user := &po.User{Uid: video.AuthorUid}
	rdb = v.db.Model(&po.User{}).Where(user).First(user)
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
	v.db.Model(&po.Video{}).Where(video).Count(&cnt)
	return cnt > 0
}

func (v *VideoDao) Insert(video *po.Video) database.DbStatus {
	rdb := v.db.Model(&po.Video{}).Model(&po.Video{}).Create(video)
	if database.IsDuplicateError(rdb.Error) {
		return database.DbExisted
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (v *VideoDao) Update(video *po.Video) database.DbStatus {
	rdb := v.db.Model(&po.Video{}).Update(video)
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
	rdb := v.db.Model(&po.Video{}).Delete(video)
	if rdb.Error != nil {
		return database.DbFailed
	} else if rdb.RowsAffected == 0 {
		return database.DbNotFound
	}
	return database.DbSuccess
}
