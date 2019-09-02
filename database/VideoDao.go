package database

import (
	. "vid/exceptions"
	. "vid/models"
)

type VideoDao struct{}

const (
	col_video_vid         = "vid"
	col_video_title       = "title"
	col_video_description = "description"
	col_video_video_url   = "video_url"
	col_video_author_uid  = "author_uid"
	col_video_upload_time = "upload_time"
)

// db 查询所有视频
//
// @return `[]Video`
func (v *VideoDao) QueryVideos() (videos []Video) {
	DB.Find(&videos)
	return videos
}

// db 查询用户视频
//
// @return `[]Video` `err`
//
// @error `UserNotExistException`
func (v *VideoDao) QueryVideosByUid(uid int) ([]Video, error) {
	var videos []Video
	var userDao UserDao
	if _, ok := userDao.QueryUserByUid(uid); !ok {
		return nil, UserNotExistException
	}
	DB.Where(col_video_author_uid+" = ?", uid).Find(&videos)
	return videos, nil
}

// db 查询 vid 视频
//
// @return `*Video` `isExist`
func (v *VideoDao) QueryVideoByVid(vid int) (*Video, bool) {
	var video Video
	DB.Where(col_video_vid+" = ?", vid).Find(&video)
	if !video.CheckValid() {
		return nil, false
	} else {
		return &video, true
	}
}
