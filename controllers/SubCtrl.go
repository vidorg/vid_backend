package controllers

import (
	"fmt"
	"net/http"
	
	. "vid/exceptions"
	. "vid/models"

	"github.com/gin-gonic/gin"
)

type SubCtrl struct{}

// POST /user/sub?up_uid&subscriber_uid
func (u *SubCtrl) SubscribeUser(c *gin.Context) {
	up_uid, ok := reqUtil.GetIntQuery(c, "up_uid")
	subscriber_uid, ok2 := reqUtil.GetIntQuery(c, "subscriber_uid")

	lostParam := make([]string, 1, 2)
	if !ok {
		lostParam = append(lostParam, "up_uid")
	}
	if !ok2 {
		lostParam = append(lostParam, "subscriber_uid")
	}

	if !ok || !ok2 {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), lostParam),
		})
		return
	}

	err := userDao.SubscribeUser(up_uid, subscriber_uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, Message{
			Message: fmt.Sprintf("User ID: %d Subscribe User ID: %d Success", subscriber_uid, up_uid),
		})
	}
}

// POST /user/unsub?up_uid&subscriber_uid
func (u *SubCtrl) UnSubscribeUser(c *gin.Context) {
	up_uid, ok := reqUtil.GetIntQuery(c, "up_uid")
	subscriber_uid, ok2 := reqUtil.GetIntQuery(c, "subscriber_uid")

	lostParam := make([]string, 1, 2)
	if !ok {
		lostParam = append(lostParam, "up_uid")
	}
	if !ok2 {
		lostParam = append(lostParam, "subscriber_uid")
	}

	if !ok || !ok2 {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), lostParam),
		})
		return
	}

	err := userDao.UnSubscribeUser(up_uid, subscriber_uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, Message{
			Message: fmt.Sprintf("User ID: %d UnSubscribe User ID: %d Success", subscriber_uid, up_uid),
		})
	}
}

// GET /user/subscriber/:uid
func (u *SubCtrl) QuerySubscriberUsers(c *gin.Context) {
	uid, ok := reqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := userDao.QuerySubscriberUsers(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// GET /user/subscriber/:uid
func (u *SubCtrl) QuerySubscribingUsers(c *gin.Context) {
	uid, ok := reqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := userDao.QuerySubscribingUsers(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}
