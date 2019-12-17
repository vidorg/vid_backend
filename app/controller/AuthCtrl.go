package controller

import (
	"net/http"
	"strconv"
	"vid/app/controller/exception"
	"vid/app/database"
	"vid/app/database/dao"
	"vid/app/middleware"
	"vid/app/model/dto"
	"vid/app/model/po"
	"vid/app/util"

	"github.com/gin-gonic/gin"
)

type authCtrl struct{}

var AuthCtrl = new(authCtrl)

// @Router 				/auth/login [POST]
// @Security 			ApiKeyAuth
// @Summary 			登录
// @Description 		用户登录
// @Param 				username formData string true 用户名
// @Param 				password formData string true 用户密码
// @Param 				expire formData integer false 登录有效期，默认一个小时
// @ID 					auth-login
// @Accept 				multipart/form-data
// @Produce 			application/json
/* @Success 200 		{
							"code": 200,
							"message": "Success",
							"data": {}
 						} */
/* @Failure 400 	 	{
							"code": 400,
							"message": "request form data exception"
 						} */
/* @Failure 401 		{
							"code": 401,
							"message": "password error"
 						} */
/* @Failure 404 		{
							"code": 404,
							"message": "user not found"
 						} */
/* @Failure 500 		{
							"code": 500,
							"message": "login failed"
 						} */
func (u *authCtrl) Login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	expireString := c.PostForm("expire")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}
	expire := util.JwtExpire
	if val, err := strconv.Atoi(expireString); err == nil {
		expire = int64(val)
	}

	passRecord := dao.PassDao.QueryByUsername(username)
	if passRecord == nil {
		c.JSON(http.StatusNotFound,
			dto.Result{}.Error(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()))
		return
	}

	if !util.PassUtil.MD5Check(password, passRecord.EncryptedPass) {
		c.JSON(http.StatusUnauthorized,
			dto.Result{}.Error(http.StatusUnauthorized).SetMessage(exception.PasswordError.Error()))
		return
	}

	token, err := util.PassUtil.GenToken(passRecord.User.Uid, expire)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.LoginError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("user", passRecord.User).PutData("token", token))
}

// POST /auth/register (Non-Auth)
func (u *authCtrl) Register(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}

	passRecord := &po.PassRecord{
		EncryptedPass: util.PassUtil.MD5Encode(password),
		User: &po.User{
			Username: username, RegisterIP: c.ClientIP(),
		},
	}
	status := dao.PassDao.Insert(passRecord)
	if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.RegisterError.Error()))
		return
	} else if status == database.DbExisted {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserNameUsedError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("user", passRecord.User))
}

// @securityDefinitions.basic BasicAuth
// POST /auth/pass (Auth)
func (u *authCtrl) ModifyPass(c *gin.Context) {
	user := middleware.GetAuthUser(c)

	password := c.PostForm("password")
	if password == "" {
		c.JSON(http.StatusBadRequest,
			dto.Result{}.Error(http.StatusBadRequest).SetMessage(exception.FormParamError.Error()))
		return
	}

	passRecord := &po.PassRecord{
		EncryptedPass: util.PassUtil.MD5Encode(password),
		User:          user,
		Uid:           user.Uid,
	}
	status := dao.PassDao.Update(passRecord)
	if status == database.DbNotFound {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UserNameUsedError.Error()))
		return
	} else if status == database.DbFailed {
		c.JSON(http.StatusInternalServerError,
			dto.Result{}.Error(http.StatusInternalServerError).SetMessage(exception.UpdatePassError.Error()))
		return
	}

	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("user", passRecord.User))
}

// GET /auth/ (Auth)
func (u *authCtrl) CurrentUser(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	c.JSON(http.StatusOK,
		dto.Result{}.Ok().PutData("user", user))
}
