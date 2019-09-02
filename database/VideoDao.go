package database

import (
	. "vid/exceptions"
	. "vid/models"
)

type videoDao struct{}

var VideoDao = new(videoDao)

const (
	col_video_vid         = "vid"
	col_video_title       = "title"
	col_video_description = "description"
	col_video_video_url   = "video_url"
	col_video_author_uid  = "author_uid"
	col_video_upload_time = "upload_time"
)

// db 查询所有视频和作者
//
// @return `[]Video`
func (v *videoDao) QueryVideos() (videos []Video) {
	DB.Find(&videos)
	for k, v := range videos {
		user, ok := UserDao.QueryUserByUid(v.AuthorUid)
		if ok {
			videos[k].Author = user
		}
	}
	return videos
}

// db 查询用户视频和作者
//
// @return `[]Video` `err`
//
// @error `UserNotExistException`
func (v *videoDao) QueryVideosByUid(uid int) ([]Video, error) {
	var videos []Video
	if _, ok := UserDao.QueryUserByUid(uid); !ok {
		return nil, UserNotExistException
	}
	DB.Where(col_video_author_uid+" = ?", uid).Find(&videos)
	user, ok := UserDao.QueryUserByUid(uid)
	if ok {
		for k, _ := range videos {
			videos[k].Author = user
		}
	}
	return videos, nil
}

// db 查询 vid 视频和作者
//
// @return `*Video` `isExist`
func (v *videoDao) QueryVideoByVid(vid int) (*Video, bool) {
	var video Video
	DB.Where(col_video_vid+" = ?", vid).Find(&video)
	if !video.CheckValid() {
		return nil, false
	} else {
		user, ok := UserDao.QueryUserByUid(video.AuthorUid)
		if ok {
			video.Author = user
		}
		return &video, true
	}
}
