package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vid/app/controller/exception"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model/dto"
)

type subCtrl struct{}

var SubCtrl = new(subCtrl)

// GET /user/:uid/subscriber?page (Non-Auth)
func (u *subCtrl) QuerySubscriberUsers(c *gin.Context) {
	uidString, _ := c.Params.Get("uid")
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

	users, count, status := dao.SubDao.QuerySubscriberUsers(uid, page)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetPage(count, page, users))
}

// GET /user/:uid/subscribing?page (Non-Auth)
func (u *subCtrl) QuerySubscribingUsers(c *gin.Context) {
	uidString, _ := c.Params.Get("uid")
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

	users, count, status := dao.SubDao.QuerySubscribingUsers(uid, page)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetPage(count, page, users))
}

// POST /user/sub?uid (Auth)
func (u *subCtrl) SubscribeUser(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	upUidString := c.Query("uid")
	upUid, err := strconv.Atoi(upUidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}

	status := dao.SubDao.SubscribeUser(user.Uid, upUid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbExtra {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.SubscribeSelfError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("me", user.Uid).PutData("up", upUid).PutData("action", "subscribe"))
}

// POST /user/unsub?uid (Auth)
func (u *subCtrl) UnSubscribeUser(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	upUidString := c.Query("uid")
	upUid, err := strconv.Atoi(upUidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}

	status := dao.SubDao.UnSubscribeUser(user.Uid, upUid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbExtra {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.SubscribeSelfError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("me", user.Uid).PutData("up", upUid).PutData("action", "unsubscribe"))
}
