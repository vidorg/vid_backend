package dao

import (
	"log"
	. "vid/app/database"
	"vid/app/model/po"
)

type videoDao struct{}

var VideoDao = new(videoDao)

func (v *videoDao) QueryAll(page int) (videos []po.Video, count int) {
	DB.Model(&po.Video{}).Count(&count)
	DB.Limit(PageSize).Offset((page - 1) * PageSize).Find(&videos)
	for idx := range videos {
		author := &po.User{}
		DB.Where(&po.User{Uid: videos[idx].AuthorUid}).Find(author)
		videos[idx].Author = author
	}
	return videos, count
}

func (v *videoDao) QueryByUid(uid int, page int) (videos []po.Video, count int, status DbStatus) {
	author := UserDao.QueryByUid(uid)
	if author == nil {
		return nil, 0, DbNotFound
	}
	video := &po.Video{AuthorUid: uid}
	DB.Where(video).Count(&count)
	DB.Limit(PageSize).Offset((page - 1) * PageSize).Where(video).Find(&videos)
	for idx := range videos {
		videos[idx].Author = author
	}
	return videos, count, DbSuccess
}

func (v *videoDao) QueryCount(uid int) (int, DbStatus) {
	if DB.NewRecord(&po.User{Uid: uid}) {
		return 0, DbNotFound
	}
	var count int
	DB.Where(&po.Video{AuthorUid: uid}).Count(&count)
	return count, DbSuccess
}

func (v *videoDao) QueryByVid(vid int) *po.Video {
	video := &po.Video{Vid: vid}
	if DB.NewRecord(video) || DB.Where(video).First(video).RecordNotFound() {
		return nil
	}
	user := &po.User{Uid: video.AuthorUid}
	if DB.Where(user).First(user).RecordNotFound() {
		return nil
	}
	video.Author = user
	return video
}

func (v *videoDao) Insert(video *po.Video) DbStatus {
	if err := DB.Create(video).Error; err != nil {
		if IsDuplicateError(err) {
			return DbExisted
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}

func (v *videoDao) Update(video *po.Video, uid int) DbStatus {
	if DB.NewRecord(video) {
		return DbNotFound
	}
	if video.AuthorUid != uid {
		return DbNotFound
	}
	if err := DB.Model(video).Update(video).Error; err != nil {
		if IsNotFoundError(err) {
			return DbNotFound
		} else {
			log.Println(err)
			return DbFailed
		}
	}
	return DbSuccess
}

func (v *videoDao) Delete(vid int, uid int) DbStatus {
	video := &po.Video{Vid: vid}
	if DB.NewRecord(video) || DB.Where(video).First(video).RecordNotFound() {
		return DbNotFound
	}
	if video.AuthorUid != uid {
		return DbNotFound
	}
	if err := DB.Model(video).Update(video).Error; err != nil {
		log.Println(err)
		return DbFailed
	}
	return DbSuccess
}
