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

// @Router 				/user/{uid}/subscriber [GET]
// @Summary 			用户粉丝
/* @Description 		查询用户所有粉丝，返回分页数据，Non-Auth

						| code | message |
						| --- | --- |
						| 400 | request route param exception |
 						| 404 | user not found | */
// @Param 				uid path integer true 所查询的用户 id
// @Param 				page query integer false 分页
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {
								"count": 2,
								"page": 1,
								"data": [
									{
										"uid": 1,
										"username": "User1",
										"sex": "male",
										"profile": "",
										"avatar_url": "",
										"birth_time": "2000-01-01",
										"authority": "admin"
									},
									{
										"uid": 2,
										"username": "User2",
										"sex": "unknown",
										"profile": "",
										"avatar_url": "",
										"birth_time": "2000-01-01",
										"authority": "normal"
									}
								]
							}
 						} */
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

// @Router 				/user/{uid}/subscribing [GET]
// @Summary 			用户关注的人
/* @Description 		查询用户所有关注，返回分页数据，Non-Auth

						| code | message |
						| --- | --- |
						| 400 | request route param exception |
 						| 404 | user not found | */
// @Param 				uid path integer true 所查询的用户 id
// @Param 				page query integer false 分页
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {
								"count": 1,
								"page": 1,
								"data": [
									{
										"uid": 1,
										"username": "User1",
										"sex": "male",
										"profile": "",
										"avatar_url": "",
										"birth_time": "2000-01-01",
										"authority": "admin"
									}
								]
							}
 						} */
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

// @Router 				/user/sub [POST]
// @Summary 			关注用户
/* @Description 		关注某一用户，Auth

						| code | message |
						| --- | --- |
						| 400 | request query param exception |
 						| 400 | subscribe oneself invalid |
						| 401 | authorization failed |
						| 401 | token has expired |
 						| 404 | user not found |
 						| 500 | subscribe failed | */
// @Param 				Authorization header string true 用户 Token
// @Param 				uid query integer true 对方用户 id
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {
								"me": 10,
								"up": 3,
								"action": "subscribe"
							}
 						} */
func (u *subCtrl) SubscribeUser(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	upUidString := c.Query("uid")
	upUid, err := strconv.Atoi(upUidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.QueryParamError.Error()))
		return
	}
	if user.Uid == upUid {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.SubscribeSelfError.Error()))
		return
	}

	status := dao.SubDao.SubscribeUser(user.Uid, upUid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.SubscribeError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("me", user.Uid).PutData("up", upUid).PutData("action", "subscribe"))
}

// @Router 				/user/unsub [POST]
// @Summary 			取消关注用户
/* @Description 		取消关注某一用户，Auth

						| code | message |
						| --- | --- |
						| 400 | request query param exception |
						| 401 | authorization failed |
						| 401 | token has expired |
 						| 404 | user not found |
 						| 500 | unsubscribe failed | */
// @Param 				Authorization header string true 用户 Token
// @Param 				uid query integer true 对方用户 id
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {
								"me": 10,
								"up": 3,
								"action": "unsubscribe"
							}
 						} */
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
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UnSubscribeError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("me", user.Uid).PutData("up", upUid).PutData("action", "unsubscribe"))
}
