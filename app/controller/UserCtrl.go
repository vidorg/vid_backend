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
	"vid/app/model/enum"
	"vid/app/model/vo"
)

type userCtrl struct{}

var UserCtrl = new(userCtrl)

// @Router 				/user/ [GET]
// @Summary 			查询所有用户
/* @Description 		管理员查询所有用户，返回分页数据，Admin

						| code | message |
						| --- | --- |
						| 401 | authorization failed |
						| 401 | token has expired |
 						| 401 | need admin authority | */
// @Param 				Authorization header string true 用户 Token
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
func (u *userCtrl) QueryAllUsers(c *gin.Context) {
	pageString := c.Query("page")
	page, err := strconv.Atoi(pageString)
	if err != nil || page == 0 {
		page = 1
	}
	users, count := dao.UserDao.QueryAll(page)
	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetPage(count, page, users))
}

// @Router 				/user/{uid} [GET]
// @Summary 			查询用户
/* @Description 		普通用户查询用户信息，Non-Auth

						| code | message |
						| --- | --- |
						| 400 | request route param exception |
 						| 404 | user not found | */
// @Param 				uid path integer true 用户 id
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {
								"user": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "unknown",
									"profile": "",
									"avatar_url": "",
									"birth_time": "2000-01-01",
									"authority": "admin"
								},
								"extra": {
									"subscribing_cnt": 1,
									"subscriber_cnt": 2,
									"video_cnt": 0,
									"playlist_cnt": 0
								}
							}
 						} */
func (u *userCtrl) QueryUser(c *gin.Context) {
	uidString := c.Param("uid")
	uid, err := strconv.Atoi(uidString)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	user := dao.UserDao.QueryByUid(uid)
	if user == nil {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	isSelfOrAdmin := middleware.GetAuthUser(c) == nil || user.Authority == enum.AuthAdmin
	extraInfo, _ := dao.UserDao.QueryUserExtraInfo(isSelfOrAdmin, user)

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("user", user).PutData("extra", extraInfo))
}

// @Router 				/user/ [PUT]
// @Summary 			更新用户
/* @Description 		更新用户信息，Auth

						| code | message |
						| --- | --- |
						| 401 | authorization failed |
						| 401 | token has expired |
 						| 404 | user not found |
 						| 500 | user update failed | */
// @Param 				Authorization header string true 用户 Token
// @Param 				username formData string false 新用户名
// @Param 				sex formData string false 新用户性别，只允许为 (male, female, unknown)
// @Param 				profile formData string false 新用户简介
// @Param 				birth_time formData string false 新用户生日，固定格式为 2000-01-01
// @Param 				phone_number formData string false 新用户电话号码
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {
								"uid": 10,
								"username": "aoihosizora",
								"sex": "male",
								"profile": "Demo Profile",
								"avatar_url": "",
								"birth_time": "2019-12-18",
								"authority": "admin"
							}
 						} */
func (u *userCtrl) UpdateUser(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	username := c.DefaultPostForm("username", user.Username)
	sex := enum.StringToSex(c.DefaultPostForm("sex", string(user.Sex)))
	profile := c.DefaultPostForm("profile", user.Profile)
	birthTime := vo.JsonDate{}.Parse(c.DefaultPostForm("birth_time", user.BirthTime.String()), user.BirthTime)
	phoneNumber := c.DefaultPostForm("phone_number", user.PhoneNumber)

	user.Username = username
	user.Sex = sex
	user.Profile = profile
	user.BirthTime = birthTime
	user.PhoneNumber = phoneNumber

	status := dao.UserDao.Update(user)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserUpdateError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().SetData(user))
}

// @Router 				/user/ [DELETE]
// @Summary 			删除用户
/* @Description 		删除用户所有信息，Auth

						| code | message |
						| --- | --- |
						| 401 | authorization failed |
						| 401 | token has expired |
 						| 404 | user not found |
 						| 404 | user delete failed | */
// @Param 				Authorization header string true 用户 Token
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "Success"
 						} */
func (u *userCtrl) DeleteUser(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	user, status := dao.UserDao.Delete(user.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserDeleteError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok())
}
