package dao
//
// import (
// 	. "vid/app/controller/exception"
// 	"vid/app/database"
// 	po2 "vid/app/model/po"
// 	"vid/app/util"
// )
//
// type videoDao struct{}
//
// var VideoDao = new(videoDao)
//
// // db 查询所有视频和作者
// //
// // @return `[]Video`
// func (v *videoDao) QueryVideos() (videos []po2.Video) {
// 	database.DB.Find(&videos)
// 	for k, _ := range videos {
// 		videos[k].ToServer()
// 		user, ok := UserDao.QueryUserByUid(videos[k].AuthorUid)
// 		if ok {
// 			videos[k].Author = user
// 		}
// 	}
// 	return videos
// }
//
// // db 查询用户视频和作者
// //
// // @return `[]Video` `err`
// //
// // @exception `UserNotFoundError`
// func (v *videoDao) QueryVideosByUid(uid int) ([]po2.Video, error) {
// 	var videos []po2.Video
// 	if _, ok := UserDao.QueryUserByUid(uid); !ok {
// 		return nil, UserNotFoundError
// 	}
// 	database.DB.Where(col_video_author_uid+" = ?", uid).Find(&videos)
// 	user, ok := UserDao.QueryUserByUid(uid)
// 	if ok {
// 		for k, _ := range videos {
// 			videos[k].ToServer()
// 			videos[k].Author = user
// 		}
// 	}
// 	return videos, nil
// }
//
// // db 查询用户发布视频数
// //
// // @return `video_cnt` `err`
// //
// // @exception `UserNotFoundError`
// func (v *videoDao) QueryUserVideoCnt(uid int) (int, error) {
// 	if _, ok := UserDao.QueryUserByUid(uid); !ok {
// 		return 0, UserNotFoundError
// 	}
// 	var videos []po2.Video
// 	database.DB.Where(col_video_author_uid+" = ?", uid).Find(&videos)
// 	return len(videos), nil
// }
//
// // db 查询 vid 视频和作者
// //
// // @return `*Video` `isExist`
// func (v *videoDao) QueryVideoByVid(vid int) (*po2.Video, bool) {
// 	var video po2.Video
// 	nf := database.DB.Where(col_video_vid+" = ?", vid).Find(&video).RecordNotFound()
// 	if nf {
// 		return nil, false
// 	} else {
// 		video.ToServer()
// 		user, ok := UserDao.QueryUserByUid(video.AuthorUid)
// 		if ok {
// 			video.Author = user
// 		}
// 		return &video, true
// 	}
// }
//
// // db 查询 video_url 视频
// //
// // @return `*Video` `isExist`
// func (v *videoDao) QueryVideoByUrl(video_url string) (*po2.Video, bool) {
// 	var video po2.Video
// 	nf := database.DB.Where(col_video_video_url+" = ?", video_url).Find(&video).RecordNotFound()
// 	if nf {
// 		return nil, false
// 	} else {
// 		video.ToServer()
// 		return &video, true
// 	}
// }
//
// // db 根据标题模糊查询视频
// //
// // @return `[]video`
// func (u *videoDao) SearchByVideoTitle(title string) (videos []po2.Video) {
// 	database.DB.Where(col_video_title+" like ?", "%"+title+"%").Find(&videos).RecordNotFound()
// 	for k, _ := range videos {
// 		videos[k].ToServer()
// 	}
// 	return videos
// }
//
// ///////////////////////////////////////////////////////////////////////////////////////////
//
// // db 创建新视频项
// //
// // @return `*video` `err`
// //
// // @exception `VideoUrlUsedError` `VideoInsertError`
// func (v *videoDao) InsertVideo(video *po2.Video) (*po2.Video, error) {
// 	// 检查同资源
// 	if _, ok := v.QueryVideoByUrl(video.VideoUrl); ok {
// 		return nil, VideoUrlUsedError
// 	}
//
// 	if video.CoverUrl == "" {
// 		video.CoverUrl = util.CmnUtil.GetDefaultVideoCoverUrl()
// 	}
//
// 	video.ToDB()
// 	// 新建
// 	database.DB.Create(video)
// 	query, ok := v.QueryVideoByVid(video.Vid)
// 	if !ok {
// 		return nil, VideoInsertError
// 	} else {
// 		query.ToServer()
// 		return query, nil
// 	}
// }
//
// // db 更新旧视频项
// //
// // @return `*video` `err`
// //
// // @exception `VideoNotFoundError` `AuthorizationError` `NotUpdateVideoError`
// func (v *videoDao) UpdateVideo(video *po2.Video, uid int) (*po2.Video, error) {
// 	old, ok := v.QueryVideoByVid(video.Vid)
// 	if !ok {
// 		return nil, VideoNotFoundError
// 	}
//
// 	// 非作者
// 	if old.AuthorUid != uid {
// 		return nil, AuthorizationError
// 	}
//
// 	video.ToDB()
//
// 	// 更新空字段
// 	if video.Title == "" {
// 		video.Title = old.Title
// 	}
// 	// if video.Description == config.AppConfig.MagicToken {
// 	// 	video.Description = old.Description
// 	// }
// 	if video.VideoUrl == "" {
// 		video.VideoUrl = old.VideoUrl
// 	}
// 	if video.CoverUrl == "" {
// 		video.CoverUrl = old.CoverUrl
// 	}
//
// 	// 检查同资源
// 	if _, ok := v.QueryVideoByUrl(video.VideoUrl); ok && video.VideoUrl != old.VideoUrl {
// 		return nil, VideoUrlUsedError
// 	}
//
// 	database.DB.Model(video).Updates(map[string]interface{}{
// 		col_video_title:       video.Title,
// 		col_video_description: video.Description,
// 		col_video_video_url:   video.VideoUrl,
// 		col_video_cover_url:   video.CoverUrl,
// 	})
// 	after, _ := v.QueryVideoByVid(video.Vid)
// 	if old.Equals(after) {
// 		after.ToServer()
// 		return after, NotUpdateVideoError
// 	} else {
// 		after.ToServer()
// 		return after, nil
// 	}
// }
//
// // db 删除视频项
// //
// // @return `*video` `err`
// //
// // @exception `VideoNotFoundError` `AuthorizationError` `VideoDeleteError`
// func (v *videoDao) DeleteVideo(vid int, uid int) (*po2.Video, error) {
// 	query, ok := v.QueryVideoByVid(vid)
// 	if !ok {
// 		return nil, VideoNotFoundError
// 	}
//
// 	// 非作者
// 	if query.AuthorUid != uid {
// 		return nil, AuthorizationError
// 	}
//
// 	if database.DB.Delete(query).RowsAffected != 1 {
// 		return nil, VideoDeleteError
// 	} else {
// 		query.ToServer()
// 		return query, nil
// 	}
// }
