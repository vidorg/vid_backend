package controllers

import (
	"fmt"
	"net/http"
	
	. "vid/exceptions"
	. "vid/models"
	. "vid/models/resp"

	"github.com/gin-gonic/gin"
)

type SubCtrl struct{}

// POST /user/sub?uid (Auth)
func (u *SubCtrl) SubscribeUser(c *gin.Context) {
	authusr, _ := c.Get("user")

	me_uid := authusr.(User).Uid
	up_uid, ok := reqUtil.GetIntQuery(c, "uid")

	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), "uid"),
		})
		return
	}

	err := userDao.SubscribeUser(me_uid, up_uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, SubResp{
			Me: me_uid,
			Up: up_uid,
			Action: "Subscribe",
		})
	}
}

// POST /user/unsub?uid (Auth)
func (u *SubCtrl) UnSubscribeUser(c *gin.Context) {
	authusr, _ := c.Get("user")

	me_uid := authusr.(User).Uid
	up_uid, ok := reqUtil.GetIntQuery(c, "uid")

	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), "uid"),
		})
		return
	}

	err := userDao.UnSubscribeUser(me_uid, up_uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, SubResp{
			Me: me_uid,
			Up: up_uid,
			Action: "UnSubscribe",
		})
	}
}

// GET /user/subscriber/:uid (Non-Auth)
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

// GET /user/subscriber/:uid (Non-Auth)
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
