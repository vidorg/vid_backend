package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vid/app/controller/exception"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model"
	"vid/app/model/dto"
	"vid/app/model/enum"
	"vid/app/model/vo"
)

type userCtrl struct{}

var UserCtrl = new(userCtrl)

// @Router 				/user?page [GET] [Auth]
// @Summary 			查询所有用户
// @Description 		管理员查询所有用户，返回分页数据，Admin
// @Param 				page query integer false "分页"
// @Accept 				multipart/form-data
// @ErrorCode			401 need admin authority
/* @Success 200 		{
							"code": 200,
							"message": "success",
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
	c.JSON(http.StatusOK, dto.Result{}.Ok().SetPage(count, page, users))
}

// @Router 				/user/{uid} [GET]
// @Summary 			查询用户
// @Description 		查询用户信息
// @Param 				uid path integer true "用户id"
// @Accept 				multipart/form-data
// @ErrorCode			400 request route param error
// @ErrorCode			404 user not found
/* @Success 200 		{
							"code": 200,
							"message": "success",
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
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RouteParamError.Error()))
		return
	}

	user := dao.UserDao.QueryByUid(uid)
	if user == nil {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}
	authUser := middleware.GetAuthUser(c)
	isSelfOrAdmin := authUser != nil && (authUser.Uid == uid || authUser.Authority == enum.AuthAdmin)
	if !isSelfOrAdmin {
		user.PhoneNumber = ""
	}
	subscribingCnt, subscriberCnt, _ := dao.SubDao.QuerySubCnt(user.Uid)
	extraInfo := &dto.UserExtraInfo{
		PhoneNumber:      user.PhoneNumber,
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       0, // TODO
		PlaylistCount:    0,
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok().PutData("user", user).PutData("extra", extraInfo))
}

// @Router 				/user/ [PUT] [Auth]
// @Summary 			更新用户
// @Description 		更新用户信息
// @Param 				username formData string true "用户名" minLength(8) maxLength(30)
// @Param 				sex formData string true "用户性别" enum(male, female, unknown)
// @Param 				profile formData string true "用户简介" minLength(0) maxLength(255)
// @Param 				birth_time formData string true "用户生日，固定格式为2000-01-01"
// @Param 				phone_number formData string true "用户手机号码"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request format error
// @ErrorCode 			400 username has been used
// @ErrorCode 			404 user not found
// @ErrorCode 			500 user update failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
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
	authUser := middleware.GetAuthUser(c)

	username, exist1 := c.GetPostForm("username")
	profile, exist2 := c.GetPostForm("profile")
	sex, exist3 := c.GetPostForm("sex")
	birthTimeStr, exist4 := c.GetPostForm("birth_time")
	birthTime, err := vo.JsonDate{}.Parse(birthTimeStr)
	phoneNumber, exist5 := c.GetPostForm("phone_number")
	if !exist1 || !exist2 || !exist3 || !exist4 || !exist5 {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}
	if !model.FormatCheck.Username(username) || !model.FormatCheck.UserProfile(profile) || err != nil {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormatError.Error()))
		return
	}

	authUser.Username = username
	authUser.Sex = enum.StringToSex(sex)
	authUser.Profile = profile
	authUser.BirthTime = birthTime
	authUser.PhoneNumber = phoneNumber

	status := dao.UserDao.Update(authUser)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbExisted {
		c.JSON(http.StatusBadRequest, dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.UserNameUsedError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserUpdateError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok().SetData(authUser))
}

// @Router 				/user/ [DELETE] [Auth]
// @Summary 			删除用户
// @Description 		删除用户所有信息
// @Accept 				multipart/form-data
// @ErrorCode 			404 user not found
// @ErrorCode 			500 user delete failed
/* @Success 200 		{
							"code": 200,
							"message": "success"
 						} */
func (u *userCtrl) DeleteUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	status := dao.UserDao.Delete(authUser.Uid)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound, dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserDeleteError.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Result{}.Ok())
}
