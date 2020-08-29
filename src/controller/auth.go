package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
	"github.com/vidorg/vid_backend/src/util"
)

func init() {
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("POST", "/v1/auth/login", "登录").
			Tags("Authorization").
			Params(goapidoc.NewBodyParam("param", "LoginParam", true, "登录参数")).
			Responses(goapidoc.NewResponse(200, "_Result<LoginDto>")),

		goapidoc.NewRoutePath("POST", "/v1/auth/register", "注册").
			Tags("Authorization").
			Params(goapidoc.NewBodyParam("param", "RegisterParam", true, "注册参数")).
			Responses(goapidoc.NewResponse(200, "_Result<UserDto>")),

		goapidoc.NewRoutePath("GET", "/v1/auth/", "当前登录用户").
			Tags("Authorization").
			Securities("Jwt").
			Responses(goapidoc.NewResponse(200, "_Result<UserDto>")),

		goapidoc.NewRoutePath("POST", "/v1/auth/logout", "注销").
			Tags("Authorization").
			Securities("Jwt").
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("PUT", "/v1/auth/password", "修改密码").
			Tags("Authorization").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("param", "PasswordParam", true, "修改密码参数")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type AuthController struct {
	config         *config.Config
	jwtService     *service.JwtService
	accountService *service.AccountService
	tokenService   *service.TokenService
}

func NewAuthController() *AuthController {
	return &AuthController{
		config:         xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:     xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		accountService: xdi.GetByNameForce(sn.SAccountService).(*service.AccountService),
		tokenService:   xdi.GetByNameForce(sn.STokenService).(*service.TokenService),
	}
}

// POST /v1/auth/login
func (a *AuthController) Login(c *gin.Context) *result.Result {
	loginParam := &param.LoginParam{}
	if err := c.ShouldBind(loginParam); err != nil {
		return result.Error(exception.RequestParamError)
	}

	account := a.accountService.QueryByUsername(loginParam.Username)
	if account == nil {
		return result.Error(exception.UserNotFoundError)
	}
	if !util.AuthUtil.CheckPassword(loginParam.Password, account.EncryptedPass) {
		return result.Error(exception.WrongPasswordError)
	}

	token, err := util.AuthUtil.GenerateToken(account.User.Uid, a.config.Jwt.Expire, a.config.Jwt)
	if err != nil {
		return result.Error(exception.LoginError)
	}
	ok := a.tokenService.Insert(token, account.Uid, a.config.Jwt.Expire)
	if !ok {
		return result.Error(exception.LoginError)
	}

	ret := dto.BuildLoginDto(account.User, token)
	return result.Ok().SetData(ret)
}

// POST /v1/auth/register
func (a *AuthController) Register(c *gin.Context) *result.Result {
	registerParam := &param.RegisterParam{}
	if err := c.ShouldBind(registerParam); err != nil {
		return result.Error(exception.WrapValidationError(err))
	}

	encrypted, err := util.AuthUtil.EncryptPassword(registerParam.Password)
	if err != nil {
		return result.Error(exception.RegisterError)
	}

	account := &po.Account{
		EncryptedPass: encrypted,
		User: &po.User{
			Username:   registerParam.Username,
			RegisterIP: c.ClientIP(),
		},
	}

	status := a.accountService.Insert(account) // cascade
	if status == xstatus.DbExisted {
		return result.Error(exception.UsernameUsedError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.RegisterError)
	}

	ret := dto.BuildUserDto(account.User)
	return result.Created().SetData(ret)
}

// GET /v1/auth
func (a *AuthController) CurrentUser(c *gin.Context) *result.Result {
	user := a.jwtService.GetContextUser(c)

	ret := dto.BuildUserDto(user)
	return result.Ok().SetData(ret)
}

// POST /v1/auth/logout
func (a *AuthController) Logout(c *gin.Context) *result.Result {
	token := a.jwtService.GetToken(c)
	ok := a.tokenService.Delete(token)
	if !ok {
		return result.Error(exception.LogoutError)
	}

	return result.Ok()
}

// PUT /v1/auth/password
func (a *AuthController) UpdatePassword(c *gin.Context) *result.Result {
	authUser := a.jwtService.GetContextUser(c)
	passParam := &param.PassParam{}
	if err := c.ShouldBind(passParam); err != nil {
		return result.Error(exception.WrapValidationError(err))
	}

	encrypted, err := util.AuthUtil.EncryptPassword(passParam.Password)
	if err != nil {
		return result.Error(exception.UpdatePassError)
	}

	account := &po.Account{
		EncryptedPass: encrypted,
		Uid:           authUser.Uid,
	}
	status := a.accountService.Update(account)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UpdatePassError)
	}
	_ = a.tokenService.DeleteAll(authUser.Uid)

	return result.Ok()
}
