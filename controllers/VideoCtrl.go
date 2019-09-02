package controllers

import (
	"net/http"
	"fmt"
	
	. "vid/database"
	. "vid/utils"
	. "vid/exceptions"
	. "vid/models/resp"

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
		c.JSON(http.StatusOK,query)
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
		user, ok := UserDao.QueryUserByUid(query.AuthorUid)
		if ok {
			query.Author = user
		}
		c.JSON(http.StatusOK,query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: VideoNotExistException.Error(),
		})
	}
}
