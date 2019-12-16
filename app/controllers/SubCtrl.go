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

type subCtrl struct{}

var SubCtrl = new(subCtrl)

// POST /user/sub?uid (Auth)
func (u *subCtrl) SubscribeUser(c *gin.Context) {
	authusr, _ := c.Get("user")

	me_uid := authusr.(po2.User).Uid
	up_uid, ok := utils.ReqUtil.GetIntQuery(c, "uid")

	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.QueryParamError.Error(), "uid"),
		})
		return
	}

	err := dao.UserDao.SubscribeUser(me_uid, up_uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, resp.SubResp{
			Me:     me_uid,
			Up:     up_uid,
			Action: "Subscribe",
		})
	}
}

// POST /user/unsub?uid (Auth)
func (u *subCtrl) UnSubscribeUser(c *gin.Context) {
	authusr, _ := c.Get("user")

	me_uid := authusr.(po2.User).Uid
	up_uid, ok := utils.ReqUtil.GetIntQuery(c, "uid")

	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.QueryParamError.Error(), "uid"),
		})
		return
	}

	err := dao.UserDao.UnSubscribeUser(me_uid, up_uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, resp.SubResp{
			Me:     me_uid,
			Up:     up_uid,
			Action: "UnSubscribe",
		})
	}
}

// GET /user/subscriber/:uid (Non-Auth)
func (u *subCtrl) QuerySubscriberUsers(c *gin.Context) {
	uid, ok := utils.ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := dao.UserDao.QuerySubscriberUsers(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// GET /user/subscriber/:uid (Non-Auth)
func (u *subCtrl) QuerySubscribingUsers(c *gin.Context) {
	uid, ok := utils.ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, resp.Message{
			Message: fmt.Sprintf(exceptions.RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := dao.UserDao.QuerySubscribingUsers(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}
