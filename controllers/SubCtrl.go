package controllers

import (
	"fmt"
	"net/http"
	
	. "vid/exceptions"
	. "vid/models"
	. "vid/database"
	. "vid/utils"
	. "vid/models/resp"

	"github.com/gin-gonic/gin"
)

type subCtrl struct{}

var SubCtrl = new(subCtrl)

// POST /user/sub?uid (Auth)
func (u *subCtrl) SubscribeUser(c *gin.Context) {
	authusr, _ := c.Get("user")

	me_uid := authusr.(User).Uid
	up_uid, ok := ReqUtil.GetIntQuery(c, "uid")

	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), "uid"),
		})
		return
	}

	err := UserDao.SubscribeUser(me_uid, up_uid)
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
func (u *subCtrl) UnSubscribeUser(c *gin.Context) {
	authusr, _ := c.Get("user")

	me_uid := authusr.(User).Uid
	up_uid, ok := ReqUtil.GetIntQuery(c, "uid")

	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), "uid"),
		})
		return
	}

	err := UserDao.UnSubscribeUser(me_uid, up_uid)
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
func (u *subCtrl) QuerySubscriberUsers(c *gin.Context) {
	uid, ok := ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := UserDao.QuerySubscriberUsers(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// GET /user/subscriber/:uid (Non-Auth)
func (u *subCtrl) QuerySubscribingUsers(c *gin.Context) {
	uid, ok := ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := UserDao.QuerySubscribingUsers(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}
