package controllers

import (
	"net/http"
	"vid/database"
	"fmt"
	. "vid/exceptions"
	. "vid/models/resp"

	"github.com/gin-gonic/gin"
)

type VideoCtrl struct{}

var videoDao = new(database.VideoDao)

// GET /video/all (Non-Auth)
func (v *VideoCtrl) GetAllVideos(c *gin.Context) {
	c.JSON(http.StatusOK, videoDao.QueryVideos())
}

// GET /video/uid/:uid (Non-Auth)
func (v *VideoCtrl) GetVideosByUid(c *gin.Context) {
	uid, ok := reqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := videoDao.QueryVideosByUid(uid)
	if err == nil {
		c.JSON(http.StatusOK,query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: err.Error(),
		})
	}
}

// GET /video/vid/:vid (Non-Auth)
func (v *VideoCtrl) GetVideoByVid(c *gin.Context) {
	vid, ok := reqUtil.GetIntParam(c.Params, "vid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "vid"),
		})
		return
	}
	query, ok := videoDao.QueryVideoByVid(vid)
	if ok {
		user, ok := userDao.QueryUserByUid(query.AuthorUid)
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
