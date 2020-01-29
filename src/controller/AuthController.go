package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/database/dao"
	"github.com/vidorg/vid_backend/src/middleware"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/util"
	"net/http"
)

type AuthController struct {
	Config     *config.ServerConfig   `di:"~"`
	JwtService *middleware.JwtService `di:"~"`
	AccountDao *dao.AccountDao        `di:"~"`
	TokenDao   *dao.TokenDao          `di:"~"`
	Mapper     *xmapper.EntityMapper  `di:"~"`
}

func NewAuthController(dic *xdi.DiContainer) *AuthController {
	ctrl := &AuthController{}
	if !dic.Inject(ctrl) {
		panic("Inject failed")
	}
	return ctrl
}

// @Router              /v1/auth/login [POST]
// @Summary             登录
// @Description         用户登录
// @Tag                 Authorization
// @Param               username formData string true "用户名"
// @Param               password formData string true "用户密码"
// @Param               expire formData integer false "登录有效期，默认为七天"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           401 password error
// @ErrorCode           404 user not found
// @ErrorCode           500 login failed
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": {
								"user": ${user},
								"token": "Bearer xxx",
								"expire": 604800
							}
 						} */
func (a *AuthController) Login(c *gin.Context) {
	loginParam := &param.LoginParam{}
	if err := c.ShouldBind(loginParam); err != nil {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.RequestParamError.Error()).JSON(c) // Login only use param error
		return
	}
	if loginParam.Expire <= 0 {
		loginParam.Expire = a.Config.JwtConfig.Expire
	}

	account := a.AccountDao.QueryByUsername(loginParam.Username)
	if account == nil {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	}

	if !util.AuthUtil.CheckPassword(loginParam.Password, account.EncryptedPass) {
		result.Result{}.Result(http.StatusUnauthorized).SetMessage(exception.PasswordError.Error()).JSON(c)
		return
	}

	token, err := util.AuthUtil.GenerateToken(account.User.Uid, loginParam.Expire, a.Config.JwtConfig)
	if err != nil {
		result.Result{}.Error().SetMessage(exception.LoginError.Error()).JSON(c)
		return
	}

	ok := a.TokenDao.Insert(token, account.Uid, loginParam.Expire)
	if !ok {
		result.Result{}.Error().SetMessage(exception.LoginError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(a.Mapper.Map(&dto.UserDto{}, account.User)).(*dto.UserDto)
	result.Result{}.Ok().
		PutData("user", retDto).
		PutData("token", token).
		PutData("expire", loginParam.Expire).JSON(c)
}

// @Router              /v1/auth/register [POST]
// @Summary             注册
// @Description         注册新用户
// @Tag                 Authorization
// @Param               username formData string true "用户名，长度在 [5, 30] 之间"
// @Param               password formData string true "用户密码，长度在 [8, 30] 之间"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           400 request format error
// @ErrorCode           500 username has been used
// @ErrorCode           500 register failed
/* @Response 200		{
							"code": 201,
							"message": "created",
							"data": ${user}
 						} */
func (a *AuthController) Register(c *gin.Context) {
	registerParam := &param.RegisterParam{}
	if err := c.ShouldBind(registerParam); err != nil {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c) // Register use wrap error
		return
	}

	encrypted, err := util.AuthUtil.EncryptPassword(registerParam.Password)
	if err != nil {
		result.Result{}.Error().SetMessage(exception.RegisterError.Error()).JSON(c)
		return
	}
	passRecord := &po.Account{
		EncryptedPass: encrypted,
		User: &po.User{
			Username:   registerParam.Username,
			RegisterIP: c.ClientIP(),
		},
	}
	status := a.AccountDao.Insert(passRecord)
	if status == database.DbExisted {
		result.Result{}.Error().SetMessage(exception.UsernameUsedError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Result{}.Error().SetMessage(exception.RegisterError.Error()).JSON(c)
		return
	}

	retDto := xcondition.First(a.Mapper.Map(&dto.UserDto{}, passRecord.User)).(*dto.UserDto)
	result.Result{}.Result(http.StatusCreated).SetData(retDto).JSON(c)
}

// @Router              /v1/auth/ [GET]
// @Security            Jwt
// @Template            Auth
// @Summary             当前登录用户
// @Description         根据认证信息，查看当前登录用户
// @Tag                 Authorization
// @Accept              multipart/form-data
/* @Response 200		{
							"code": 200,
							"message": "success",
							"data": ${user}
 						} */
func (a *AuthController) CurrentUser(c *gin.Context) {
	authUser := a.JwtService.GetAuthUser(c)
	retDto := xcondition.First(a.Mapper.Map(&dto.UserDto{}, authUser)).(*dto.UserDto)
	result.Result{}.Ok().SetData(retDto).JSON(c)
}

// @Router              /v1/auth/logout [POST]
// @Security            Jwt
// @Template            Auth
// @Summary             注销
// @Description         用户注销
// @Tag                 Authorization
// @Accept              multipart/form-data
// @ErrorCode           500 logout failed
/* @Response 200		{
							"code": 200,
							"message": "success"
 						} */
func (a *AuthController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	// only delete current token
	ok := a.TokenDao.Delete(authHeader)
	if !ok {
		result.Result{}.Error().SetMessage(exception.LogoutError.Error()).JSON(c)
		return
	}

	result.Result{}.Ok().JSON(c)
}

// @Router              /v1/auth/password [PUT] [Auth]
// @Security            Jwt
// @Template            Auth
// @Summary             修改密码
// @Description         用户修改密码
// @Tag                 Authorization
// @Param               password formData string true "用户密码，长度在 [8, 30] 之间"
// @Accept              multipart/form-data
// @ErrorCode           400 request param error
// @ErrorCode           400 request format error
// @ErrorCode           404 user not found
// @ErrorCode           500 update password failed
/* @Response 200		{
							"code": 200,
							"message": "success"
 						} */
func (a *AuthController) UpdatePassword(c *gin.Context) {
	authUser := a.JwtService.GetAuthUser(c)
	passParam := &param.PassParam{}
	if err := c.ShouldBind(passParam); err != nil {
		result.Result{}.Result(http.StatusBadRequest).SetMessage(exception.WrapValidationError(err).Error()).JSON(c)
		return
	}

	encrypted, err := util.AuthUtil.EncryptPassword(passParam.Password)
	if err != nil {
		result.Result{}.Error().SetMessage(exception.UpdatePassError.Error()).JSON(c)
		return
	}
	passRecord := &po.Account{
		EncryptedPass: encrypted,
		Uid:           authUser.Uid,
	}
	status := a.AccountDao.Update(passRecord)
	if status == database.DbNotFound {
		result.Result{}.Result(http.StatusNotFound).SetMessage(exception.UserNotFoundError.Error()).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Result{}.Error().SetMessage(exception.UpdatePassError.Error()).JSON(c)
		return
	}
	// Delete all token but ignore result
	_ = a.TokenDao.DeleteAll(authUser.Uid)

	result.Result{}.Ok().JSON(c)
}
