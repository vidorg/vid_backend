package controller

import (
	"net/http"
	"vid/app/controller/exception"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model"
	"vid/app/model/dto"
	"vid/app/model/dto/common"
	"vid/app/model/dto/param"
	"vid/app/model/enum"
	"vid/app/model/po"
	"vid/app/util"

	"github.com/gin-gonic/gin"
)

type authController struct{}

var AuthController = new(authController)

// @Router 				/auth/login [POST]
// @Summary 			登录
// @Description 		用户登录
// @Tag					Authorization
// @Param 				username formData string true "用户名"
// @Param 				password formData string true "用户密码"
// @Param 				expire formData integer false "登录有效期，默认为七天"
// @Accept 				multipart/form-data
// @ErrorCode 			400 request param error
// @ErrorCode 			401 password error
// @ErrorCode 			404 user not found
// @ErrorCode 			500 login failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"user": {
									"uid": 10,
									"username": "aoihosizora",
									"sex": "unknown",
									"profile": "",
									"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
									"birth_time": "2000-01-01",
									"authority": "normal"
									"phone_number": "13512345678"
								},
								"token": "Bearer xxx",
								"expire": 604800
							}
 						} */
func (u *authController) Login(c *gin.Context) {
	passwordParam := param.PasswordParam{}
	if err := c.ShouldBind(&passwordParam); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	if passwordParam.Expire <= 0 {
		passwordParam.Expire = middleware.JwtExpire
	}

	passRecord := dao.PassDao.QueryByUsername(passwordParam.Username)
	if passRecord == nil {
		c.JSON(http.StatusNotFound, common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	if !util.AuthUtil.MD5Check(passwordParam.Password, passRecord.EncryptedPass) {
		c.JSON(http.StatusUnauthorized, common.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.PasswordError.Error()))
		return
	}

	token, err := util.AuthUtil.GenerateToken(passRecord.User.Uid, passwordParam.Expire)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.LoginError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().PutData("user", dto.UserDto{}.FromPo(passRecord.User, enum.DtoOptionAll)).
		PutData("token", token).PutData("expire", passwordParam.Expire))
}

// @Router 				/auth/register [POST]
// @Summary 			注册
// @Description 		用户注册
// @Tag					Authorization
// @Param 				username formData string true "用户名" minLength(5) maxLength(30)
// @Param 				password formData string true "用户密码" minLength(8) maxLength(30)
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			400 request format error
// @ErrorCode			500 username has been used
// @ErrorCode			500 register failed
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"uid": 10,
								"username": "aoihosizora",
								"sex": "unknown",
								"profile": "",
								"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
								"birth_time": "2000-01-01",
								"authority": "normal"
								"phone_number": "13512345678"
							}
 						} */
func (u *authController) Register(c *gin.Context) {
	passwordParam := param.PasswordParam{}
	if err := c.ShouldBind(&passwordParam); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	if !model.FormatCheck.Username(passwordParam.Username) || !model.FormatCheck.Password(passwordParam.Password) {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestFormatError.Error()))
		return
	}

	passRecord := &po.Password{
		EncryptedPass: util.AuthUtil.MD5Encode(passwordParam.Password),
		User: &po.User{
			Username:   passwordParam.Username,
			RegisterIP: c.ClientIP(),
		},
	}
	status := dao.PassDao.Insert(passRecord)
	if status == database.DbExisted {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserNameUsedError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError, common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.RegisterError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok().SetData(dto.UserDto{}.FromPo(passRecord.User, enum.DtoOptionAll)))
}

// @Router 				/auth/password [PUT] [Auth]
// @Summary 			修改密码
// @Description 		用户修改密码
// @Tag					Authorization
// @Param 				password formData string true "用户密码" minLength(8) maxLength(30)
// @Accept 				multipart/form-data
// @ErrorCode			400 request param error
// @ErrorCode			400 request format error
// @ErrorCode			404 user not found
// @ErrorCode			500 update password failed
/* @Success 200 		{
							"code": 200,
							"message": "success"
 						} */
func (u *authController) ModifyPassword(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)

	password, exist := c.GetPostForm("password")
	if !exist {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()))
		return
	}
	if !model.FormatCheck.Password(password) {
		c.JSON(http.StatusBadRequest, common.Result{}.Error(http.StatusBadRequest).SetMessage(exception.RequestFormatError.Error()))
		return
	}

	passRecord := &po.Password{
		EncryptedPass: util.AuthUtil.MD5Encode(password),
		User:          authUser,
		Uid:           authUser.Uid,
	}
	status := dao.PassDao.Update(passRecord)
	if status == database.DbNotFound {
		c.JSON(http.StatusNotFound,
			common.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			common.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UpdatePassError.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Result{}.Ok())
}

// @Router 				/auth/ [GET] [Auth]
// @Summary 			查看当前登录用户
// @Description 		根据认证令牌，查看当前登录用户
// @Tag					Authorization
// @Accept 				multipart/form-data
/* @Success 200 		{
							"code": 200,
							"message": "success",
							"data": {
								"uid": 10,
								"username": "aoihosizora",
								"sex": "unknown",
								"profile": "",
								"avatar_url": "http://localhost:3344/raw/image/default/avatar.jpg",
								"birth_time": "2000-01-01",
								"authority": "normal"
								"phone_number": "13512345678"
							}
 						} */
func (u *authController) CurrentUser(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	c.JSON(http.StatusOK, common.Result{}.Ok().SetData(dto.UserDto{}.FromPo(authUser, enum.DtoOptionAll)))
}
