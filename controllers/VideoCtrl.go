package controllers

import (
	"fmt"
	"net/http"

	. "vid/database"
	. "vid/exceptions"
	. "vid/models/resp"
	. "vid/models"
	. "vid/utils"

	"github.com/gin-gonic/gin"
)

type videoCtrl struct{}

var VideoCtrl = new(videoCtrl)

// GET /video/all (Non-Auth)
func (v *videoCtrl) GetAllVideos(c *gin.Context) {
	c.JSON(http.StatusOK, VideoDao.QueryVideos())
}

// GET /video/uid/:uid (Non-Auth)
func (v *videoCtrl) GetVideosByUid(c *gin.Context) {
	uid, ok := ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := VideoDao.QueryVideosByUid(uid)
	if err == nil {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: err.Error(),
		})
	}
}

// GET /video/vid/:vid (Non-Auth)
func (v *videoCtrl) GetVideoByVid(c *gin.Context) {
	vid, ok := ReqUtil.GetIntParam(c.Params, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "vid"),
		})
		return
	}
	query, ok := VideoDao.QueryVideoByVid(vid)
	if ok {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: VideoNotExistException.Error(),
		})
	}
}

// POST /video/new (Auth)
func (v *videoCtrl) UploadNewVideo(c *gin.Context) {
	body := ReqUtil.GetBody(c.Request.Body)
	var video Video
	if !video.Unmarshal(body, true) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid
	video.AuthorUid = uid

	query, err := VideoDao.InsertVideo(&video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// POST /video/update (Auth)
func (v *videoCtrl) UpdateVideoInfo(c *gin.Context) {
	body := ReqUtil.GetBody(c.Request.Body)
	var video Video
	if !video.Unmarshal(body, false) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}
	
	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid
	video.AuthorUid = uid

	query, err := VideoDao.UpdateVideo(&video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /video/delete?vid (Auth)
func (v *videoCtrl) DeleteVideo(c *gin.Context) {
	vid, ok := ReqUtil.GetIntQuery(c, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "vid"),
		})
		return
	}

	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid

	query, err := VideoDao.DeleteVideo(vid, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}