package controller

import (
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/ahlib/xdi"
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
func (a *AuthController) Login(c *gin.Context) {
	loginParam := &param.LoginParam{}
	if err := c.ShouldBind(loginParam); err != nil {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	account := a.accountService.QueryByUsername(loginParam.Username)
	if account == nil {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}
	if !util.AuthUtil.CheckPassword(loginParam.Password, account.EncryptedPass) {
		result.Error(exception.WrongPasswordError).JSON(c)
		return
	}

	token, err := util.AuthUtil.GenerateToken(account.User.Uid, a.config.Jwt.Expire, a.config.Jwt)
	if err != nil {
		result.Error(exception.LoginError).JSON(c)
		return
	}
	ok := a.tokenService.Insert(token, account.Uid, a.config.Jwt.Expire)
	if !ok {
		result.Error(exception.LoginError).JSON(c)
		return
	}

	ret := dto.BuildLoginDto(account.User, token)
	result.Ok().SetData(ret).JSON(c)
}

// POST /v1/auth/register
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

	account := &po.Account{
		EncryptedPass: encrypted,
		User: &po.User{
			Username:   registerParam.Username,
			RegisterIP: c.ClientIP(),
		},
	}

	status := a.accountService.Insert(account) // cascade
	if status == xstatus.DbExisted {
		result.Error(exception.UsernameUsedError).JSON(c)
		return
	} else if status == xstatus.DbFailed {
		result.Error(exception.RegisterError).JSON(c)
		return
	}

	ret := dto.BuildUserDto(account.User)
	result.Created().SetData(ret).JSON(c)
}

// GET /v1/auth
func (a *AuthController) CurrentUser(c *gin.Context) {
	user := a.jwtService.GetContextUser(c)

	ret := dto.BuildUserDto(user)
	result.Ok().SetData(ret).JSON(c)
}

// POST /v1/auth/logout
func (a *AuthController) Logout(c *gin.Context) {
	token := a.jwtService.GetToken(c)
	ok := a.tokenService.Delete(token)
	if !ok {
		result.Error(exception.LogoutError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// PUT /v1/auth/password
func (a *AuthController) UpdatePassword(c *gin.Context) {
	authUser := a.jwtService.GetContextUser(c)
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

	account := &po.Account{
		EncryptedPass: encrypted,
		Uid:           authUser.Uid,
	}
	status := a.accountService.Update(account)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == xstatus.DbFailed {
		result.Error(exception.UpdatePassError).JSON(c)
		return
	}
	_ = a.tokenService.DeleteAll(authUser.Uid)

	result.Ok().JSON(c)
}
