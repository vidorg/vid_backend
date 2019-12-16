package controllers

import (
	"fmt"
	"net/http"
	"vid/app/controllers/exceptions"
	"vid/app/database/dao"
	po2 "vid/app/models/po"
	"vid/app/models/resp"
	"vid/app/utils"

	"github.com/gin-gonic/gin"
)

type videoCtrl struct{}

var VideoCtrl = new(videoCtrl)

// GET /video/all (Auth) (Admin)
func (v *videoCtrl) GetAllVideos(c *gin.Context) {
	authusr, _ := c.Get("user")
	if authusr.(po2.User).Authority != po2.AuthAdmin {
		c.JSON(http.StatusUnauthorized, resp.Message{
			Message: exceptions.NeedAdminException.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dao.VideoDao.QueryVideos())
}

// GET /video/uid/:uid (Non-Auth)
func (v *videoCtrl) GetVideosByUid(c *gin.Context) {
	uid, ok := utils.ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := dao.VideoDao.QueryVideosByUid(uid)
	if err == nil {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, resp.Message{
			Message: err.Error(),
		})
	}
}

// GET /video/vid/:vid (Non-Auth)
func (v *videoCtrl) GetVideoByVid(c *gin.Context) {
	vid, ok := utils.ReqUtil.GetIntParam(c.Params, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "vid"),
		})
		return
	}
	query, ok := dao.VideoDao.QueryVideoByVid(vid)
	if ok {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, resp.Message{
			Message: exceptions.VideoNotExistException.Error(),
		})
	}
}

// POST /video/new (Auth)
func (v *videoCtrl) UploadNewVideo(c *gin.Context) {
	body := utils.ReqUtil.GetBody(c.Request.Body)
	var video po2.Video
	if !video.Unmarshal(body, true) {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exceptions.RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")
	uid := authusr.(po2.User).Uid
	video.AuthorUid = uid

	query, err := dao.VideoDao.InsertVideo(&video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// PUT /video/update (Auth)
func (v *videoCtrl) UpdateVideoInfo(c *gin.Context) {
	body := utils.ReqUtil.GetBody(c.Request.Body)
	var video po2.Video
	if !video.Unmarshal(body, false) {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: exceptions.RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")

	query, err := dao.VideoDao.UpdateVideo(&video, authusr.(po2.User).Uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /video/delete?vid (Auth)
func (v *videoCtrl) DeleteVideo(c *gin.Context) {
	vid, ok := utils.ReqUtil.GetIntQuery(c, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "vid"),
		})
		return
	}

	authusr, _ := c.Get("user")

	query, err := dao.VideoDao.DeleteVideo(vid, authusr.(po2.User).Uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}
