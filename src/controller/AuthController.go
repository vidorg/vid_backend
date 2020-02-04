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
// @Template            ParamA
// @Summary             登录
// @Tag                 Authorization
// @Param               param body #LoginParam true false "请求参数"
// @ResponseDesc 401    "password error"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "login failed"
// @ResponseModel 200   #Result<LoginDto>
// @ResponseEx 200      ${resp_login}
func (a *AuthController) Login(c *gin.Context) {
	loginParam := &param.LoginParam{}
	if err := c.ShouldBind(loginParam); err != nil {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	if loginParam.Expire <= 0 {
		loginParam.Expire = a.Config.JwtConfig.Expire
	}

	account := a.AccountDao.QueryByUsername(loginParam.Username)
	if account == nil {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}
	if !util.AuthUtil.CheckPassword(loginParam.Password, account.EncryptedPass) {
		result.Error(exception.PasswordError).JSON(c)
		return
	}

	token, err := util.AuthUtil.GenerateToken(account.User.Uid, loginParam.Expire, a.Config.JwtConfig)
	if err != nil {
		result.Error(exception.LoginError).JSON(c)
		return
	}
	ok := a.TokenDao.Insert(token, account.Uid, loginParam.Expire)
	if !ok {
		result.Error(exception.LoginError).JSON(c)
		return
	}

	retDto := xcondition.First(a.Mapper.Map(&dto.UserDto{}, account.User)).(*dto.UserDto)
	result.Ok().
		PutData("user", retDto).
		PutData("token", token).
		PutData("expire", loginParam.Expire).JSON(c)
}

// @Router              /v1/auth/register [POST]
// @Template            Param
// @Summary             注册
// @Tag                 Authorization
// @Param               param body #RegisterParam true false "请求参数"
// @ResponseDesc 400    "username has been used"
// @ResponseDesc 500    "register failed"
// @ResponseModel 201   #Result<UserDto>
// @Response 201        ${resp_register}
func (a *AuthController) Register(c *gin.Context) {
	registerParam := &param.RegisterParam{}
	if err := c.ShouldBind(registerParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	encrypted, err := util.AuthUtil.EncryptPassword(registerParam.Password)
	if err != nil {
		result.Error(exception.RegisterError).JSON(c)
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
		result.Error(exception.UsernameUsedError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.RegisterError).JSON(c)
		return
	}

	retDto := xcondition.First(a.Mapper.Map(&dto.UserDto{}, passRecord.User)).(*dto.UserDto)
	result.Status(http.StatusCreated).SetData(retDto).JSON(c)
}

// @Router              /v1/auth [GET]
// @Security            Jwt
// @Template            Auth
// @Summary             当前登录用户
// @Tag                 Authorization
// @ResponseModel 200   #Result<UserDto>
// @ResponseEx 200      ${resp_user}
func (a *AuthController) CurrentUser(c *gin.Context) {
	authUser := a.JwtService.GetAuthUser(c)
	retDto := xcondition.First(a.Mapper.Map(&dto.UserDto{}, authUser)).(*dto.UserDto)
	result.Ok().SetData(retDto).JSON(c)
}

// @Router              /v1/auth/logout [POST]
// @Security            Jwt
// @Template            Auth
// @Summary             注销
// @Tag                 Authorization
// @ResponseDesc 500    "logout failed"
// @ResponseModel 200   #Result
// @ResponseEx 200      ${resp_success}
func (a *AuthController) Logout(c *gin.Context) {
	authHeader := a.JwtService.GetAuthToken(c)
	ok := a.TokenDao.Delete(authHeader)
	if !ok {
		result.Error(exception.LogoutError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// @Router              /v1/auth/password [PUT]
// @Security            Jwt
// @Template            Auth Param
// @Summary             修改密码
// @Tag                 Authorization
// @Param               param body #PassParam true false "请求参数"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "update password failed"
// @ResponseModel 200   #Result
// @ResponseEx 200      ${resp_success}
func (a *AuthController) UpdatePassword(c *gin.Context) {
	authUser := a.JwtService.GetAuthUser(c)
	passParam := &param.PassParam{}
	if err := c.ShouldBind(passParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	encrypted, err := util.AuthUtil.EncryptPassword(passParam.Password)
	if err != nil {
		result.Error(exception.UpdatePassError).JSON(c)
		return
	}
	passRecord := &po.Account{
		EncryptedPass: encrypted,
		Uid:           authUser.Uid,
	}
	status := a.AccountDao.Update(passRecord)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.UpdatePassError).JSON(c)
		return
	}
	_ = a.TokenDao.DeleteAll(authUser.Uid)

	result.Ok().JSON(c)
}
