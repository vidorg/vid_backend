package controller

import (
	"github.com/gin-gonic/gin"
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
)

type videoCtrl struct{}

var VideoCtrl = new(videoCtrl)

// @Router 				/video?page [GET]
// @Summary 			查询所有视频
/* @Description 		管理员查询所有视频，返回分页数据，Admin

						| code | message |
						| --- | --- |
						| 401 | authorization failed |
						| 401 | token has expired |
 						| 401 | need admin authority | */
// @Param 				Authorization header string true 用户 Token
// @Param 				page query integer false 分页
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {}
 						} */
func (v *videoCtrl) QueryAllVideos(c *gin.Context) {
	pageString := c.Query("page")
	page, err := strconv.Atoi(pageString)
	if err != nil || page == 0 {
		page = 1
	}
	users, count := dao.VideoDao.QueryAll(page)
	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetPage(count, page, users))
}

// @Router 				/video/uid/{uid}?page [GET]
// @Summary 			查询用户视频
/* @Description 		查询作者为用户的所有视频，返回分页数据，Non-Auth

						| code | message |
						| --- | --- |
						| 400 | request route param error |
 						| 404 | user not found | */
// @Param 				uid path integer true 用户 id
// @Param 				page query integer false 分页
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {}
 						} */
func (v *videoCtrl) QueryVideosByUid(c *gin.Context) {
	uidString := c.Param("uid")
	uid, err := strconv.Atoi(uidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}
	pageString := c.Query("page")
	page, err := strconv.Atoi(pageString)
	if err != nil || page == 0 {
		page = 1
	}

	users, count, status := dao.VideoDao.QueryByUid(uid, page)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetPage(count, page, users))
}

// @Router 				/video/vid/{vid} [GET]
// @Summary 			查询视频
/* @Description 		查询视频信息，Non-Auth

						| code | message |
						| --- | --- |
						| 400 | request route param error |
 						| 404 | video not found | */
// @Param 				vid path integer true 视频 id
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {}
 						} */
func (v *videoCtrl) QueryVideoByVid(c *gin.Context) {
	vidString := c.Param("vid")
	vid, err := strconv.Atoi(vidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	video := dao.VideoDao.QueryByVid(vid)
	if video == nil {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	}
	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetData(video))
}

// @Router 				/video/ [POST]
// @Summary 			新建视频
/* @Description 		新建用户视频，Auth

						| code | message |
						| --- | --- |
						| 400 | request form data error |
						| 401 | authorization failed |
 						| 401 | token has expired |
 						| 500 | video existed failed |
 						| 401 | video insert failed | */
// @Param 				Authorization header string true 用户 Token
// @Param 				title formData string true 视频标题
// @Param 				description formData string true 视频简介
// @Param 				video_url formData string true 视频资源链接
// @Param 				cover_url formData string true 视频封面链接
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {}
 						} */
func (v *videoCtrl) InsertVideo(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	title := c.PostForm("title")
	description, exist := c.GetPostForm("description")
	url := c.PostForm("video_url")
	cover := c.PostForm("cover_url")

	if title == "" || !exist || url == "" || cover == "" {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}

	video := &po.Video{
		Title: title,
		Description: description,
		VideoUrl: url,
		CoverUrl: cover,
		UploadTime: vo.JsonDate(time.Now()),
		AuthorUid: user.Uid,
	}

	status := dao.VideoDao.Insert(video)
	if status == database.DbExisted {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoExistError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoInsertError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetData(video))
}

// @Router 				/video/{vid} [POST]
// @Summary 			更新视频
/* @Description 		更新用户视频信息，Auth

						| code | message |
						| --- | --- |
						| 400 | request route param error |
						| 400 | request format error |
						| 401 | authorization failed |
 						| 401 | token has expired |
 						| 404 | video not found |
 						| 404 | video update failed | */
// @Param 				Authorization header string true 用户 Token
// @Param 				vid path string true 更新视频 id
// @Param 				title formData string false 视频新标题
// @Param 				description formData string false 视频新简介
// @Param 				video_url formData string false 视频新资源链接
// @Param 				cover_url formData string false 视频新封面链接
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {}
 						} */
func (v *videoCtrl) UpdateVideo(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	vidString := c.Param("vid")
	vid, err := strconv.Atoi(vidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	video := dao.VideoDao.QueryByVid(vid)
	title := c.DefaultPostForm("title", video.Title)
	description := c.DefaultPostForm("description", video.Description)
	if !model.FormatCheck.VideoTitle(title) || !model.FormatCheck.VideoDesc(description) {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormatError.Error()))
		return
	}
	video.Title = title
	video.Description = description
	video.VideoUrl = c.DefaultPostForm("video_url", video.VideoUrl)
	video.CoverUrl = c.DefaultPostForm("cover_url", video.CoverUrl)
	// TODO VideoUrl CoverUrl

	status := dao.VideoDao.Update(video, user.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoUpdateError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetData(video))
}

// @Router 				/video/{vid} [DELETE]
// @Summary 			删除视频
/* @Description 		删除用户视频，Auth

						| code | message |
						| --- | --- |
						| 400 | request route param error |
						| 401 | authorization failed |
 						| 401 | token has expired |
 						| 404 | video not found |
 						| 500 | video delete failed | */
// @Param 				Authorization header string true 用户 Token
// @Param 				vid path string true 删除视频 id
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {}
 						} */
func (v *videoCtrl) DeleteVideo(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	vidString := c.Param("vid")
	vid, err := strconv.Atoi(vidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	status := dao.VideoDao.Delete(vid, user.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.VideoNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.VideoDeleteError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok())
}
