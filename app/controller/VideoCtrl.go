package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shomali11/util/xconditions"
	"net/http"
	"strconv"
	"time"
	"vid/app/controller/exception"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model"
	"vid/app/model/dto"
	"vid/app/model/po"
	"vid/app/model/vo"
	"vid/app/util"
)

type videoCtrl struct{}

var VideoCtrl = new(videoCtrl)

// @Router 				/video?page [GET] [Auth]
// @Summary 			查询所有视频
// @Description 		管理员查询所有视频，返回分页数据，Admin
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request query param error
// @ErrorCode 			401 need admin authority
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"vid": 10,
										"title": "VideoTitle",
										"description": "VideoDesc",
										"video_url": "",
										"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
										"upload_time": "2019-12-26",
										"author": {
											"uid": 10,
											"username": "aoihosizora",
											"sex": "male",
											"profile": "Demo Profile",
											"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
											"birth_time": "2019-12-26",
											"authority": "admin"
										}
									}
								]
							}
 						} */
func (v *videoCtrl) QueryAllVideos(c *gin.Context) {
	pageString := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page == 0, 1, page).(int)

	videos, count := dao.VideoDao.QueryAll(page)
	c.JSON(http.StatusOK, dto.Result{}.Ok().AddConverter(po.User{}.UrlConverter()).AddConverter(po.Video{}.UrlConverter()).SetPage(count, page, videos))
}

// @Router 				/video/uid/{uid}?page [GET]
// @Summary 			查询用户视频
// @Description 		查询作者为用户的所有视频，返回分页数据
// @Param 				uid path integer true "用户id"
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			400 request route param error
// @ErrorCode			400 request query param error
// @ErrorCode			404 user not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"vid": 10,
										"title": "VideoTitle",
										"description": "VideoDesc",
										"video_url": "",
										"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
										"upload_time": "2019-12-26",
										"author": {
											"uid": 10,
											"username": "aoihosizora",
											"sex": "male",
											"profile": "Demo Profile",
											"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
											"birth_time": "2019-12-26",
											"authority": "admin"
										}
									}
								]
							}
 						} */
func (v *videoCtrl) QueryVideosByUid(c *gin.Context) {
	uidString := c.Param("uid")
	uid, err := strconv.Atoi(uidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}
	pageString := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}
	page = xconditions.IfThenElse(page == 0, 1, page).(int)

	videos, count, status := dao.VideoDao.QueryByUid(uid, page)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok().AddConverter(po.User{}.UrlConverter()).AddConverter(po.Video{}.UrlConverter()).SetPage(count, page, videos))
}

// @Router 				/video/vid/{vid} [GET]
// @Summary 			查询视频
// @Description 		查询视频信息
// @Param 				vid path integer true "视频id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request route param error
// @ErrorCode			404 video not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"vid": 10,
								"title": "VideoTitle",
								"description": "VideoDesc",
								"video_url": "",
								"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
								"upload_time": "2019-12-26",
								"author": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "male",
									"profile": "Demo Profile",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2019-12-26",
									"authority": "admin"
								}
							}
 						} */
func (v *videoCtrl) QueryVideoByVid(c *gin.Context) {
	vidString := c.Param("vid")
	vid, err := strconv.Atoi(vidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	video := dao.VideoDao.QueryByVid(vid)
	if video == nil {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok().AddConverter(po.User{}.UrlConverter()).AddConverter(po.Video{}.UrlConverter()).SetData(video))
}

// @Router 				/video/ [POST] [Auth]
// @Summary 			新建视频
// @Description 		新建用户视频
// @Param 				title formData string true "视频标题" minLength(5) maxLength(100)
// @Param 				description formData string true "视频简介" minLength(0) maxLength(255)
// @Param 				video_url formData string true "视频资源链接"
// @Param 				cover formData file false "视频封面"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request form data error
// @ErrorCode 			400 request format error
// @ErrorCode 			400 request body too large
// @ErrorCode 			400 image type not supported
// @ErrorCode 			400 video resource has been used
// @ErrorCode 			500 image save failed
// @ErrorCode 			500 video insert failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"vid": 10,
								"title": "VideoTitle",
								"description": "VideoDesc",
								"video_url": "",
								"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
								"upload_time": "2019-12-26",
								"author": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "male",
									"profile": "Demo Profile",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2019-12-26",
									"authority": "admin"
								}
							}
 						} */
func (v *videoCtrl) InsertVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	title, exist1 := c.GetPostForm("title")
	description, exist2 := c.GetPostForm("description")
	videoUrl, exist3 := c.GetPostForm("video_url")
	if !exist1 || !exist2 || !exist3 {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}
	if !model.FormatCheck.VideoTitle(title) || !model.FormatCheck.VideoDesc(description) {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormatError.Error()))
		return
	}
	coverUrl := ""
	coverFile, coverHeader, err := c.Request.FormFile("cover")
	if err == nil && coverFile != nil {
		supported, ext := util.ImageUtil.CheckImageExt(coverHeader.Filename)
		if !supported {
			c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()))
			return
		}
		filename := fmt.Sprintf("cover_%s.jpg", util.CommonUtil.CurrentTimeUuid())
		savePath := fmt.Sprintf("./usr/image/%d/%s", authUser.Uid, filename)
		if err := util.ImageUtil.SaveAsJpg(coverFile, ext, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.ImageSaveError.Error()))
			return
		}
		coverUrl = filename
	}

	video := &po.Video{
		Title:       title,
		Description: description,
		CoverUrl:    coverUrl,
		VideoUrl:    videoUrl, // TODO
		UploadTime:  vo.JsonDate(time.Now()),
		AuthorUid:   authUser.Uid,
		Author:      authUser,
	}
	status := dao.VideoDao.Insert(video)
	if status == database.DbExisted {
		c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.VideoExistError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoInsertError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok().AddConverter(po.User{}.UrlConverter()).AddConverter(po.Video{}.UrlConverter()).SetData(video))
}

// @Router 				/video/{vid} [POST] [Auth]
// @Summary 			更新视频
// @Description 		更新用户视频信息
// @Param 				vid path string true "更新视频id"
// @Param 				title formData string false "视频标题" minLength(5) maxLength(100)
// @Param 				description formData string false "视频简介" minLength(0) maxLength(255)
// @Param 				cover formData file false "视频封面"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request route param error
// @ErrorCode 			400 request form data error
// @ErrorCode 			400 request format error
// @ErrorCode 			400 request body too large
// @ErrorCode 			400 image type not supported
// @ErrorCode 			404 video not found
// @ErrorCode 			500 image save failed
// @ErrorCode 			500 video update failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"vid": 10,
								"title": "VideoTitle",
								"description": "VideoDesc",
								"video_url": "",
								"cover_url": "http://localhost:3344/raw/image/default/cover.jpg",
								"upload_time": "2019-12-26",
								"author": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "male",
									"profile": "Demo Profile",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2019-12-26",
									"authority": "admin"
								}
							}
 						} */
func (v *videoCtrl) UpdateVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	vidString := c.Param("vid")
	vid, err := strconv.Atoi(vidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}
	title, exist1 := c.GetPostForm("title")
	description, exist2 := c.GetPostForm("description")
	if !exist1 || !exist2 {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}
	if !model.FormatCheck.VideoTitle(title) || !model.FormatCheck.VideoDesc(description) {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormatError.Error()))
		return
	}
	coverUrl := ""
	coverFile, coverHeader, err := c.Request.FormFile("cover")
	if err == nil && coverFile != nil {
		supported, ext := util.ImageUtil.CheckImageExt(coverHeader.Filename)
		if !supported {
			c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.ImageNotSupportedError.Error()))
			return
		}
		filename := fmt.Sprintf("cover_%s.jpg", util.CommonUtil.CurrentTimeUuid())
		savePath := fmt.Sprintf("./usr/image/%d/%s", authUser.Uid, filename)
		if err := util.ImageUtil.SaveAsJpg(coverFile, ext, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.ImageSaveError.Error()))
			return
		}
		coverUrl = filename
	}

	video := dao.VideoDao.QueryByVid(vid)
	video.Title = title
	video.Description = description
	video.CoverUrl = xconditions.IfThenElse(coverUrl == "", video.CoverUrl, coverUrl).(string)

	status := dao.VideoDao.Update(video, authUser.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoUpdateError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok().AddConverter(po.User{}.UrlConverter()).AddConverter(po.Video{}.UrlConverter()).SetData(video))
}

// @Router 				/video/{vid} [DELETE] [Auth]
// @Summary 			删除视频
// @Description 		删除用户视频
// @Param 				vid path string true "删除视频id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request route param error
// @ErrorCode			404 video not found
// @ErrorCode			500 video delete failed
/* @Success 200 		{
							"code": 200,
							"message": "success"
 						} */
func (v *videoCtrl) DeleteVideo(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	vidString := c.Param("vid")
	vid, err := strconv.Atoi(vidString)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	status := dao.VideoDao.Delete(vid, authUser.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoDeleteError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok())
}
