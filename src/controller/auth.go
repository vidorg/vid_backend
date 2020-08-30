package controller

import (
	"github.com/Aoi-hosizora/ahlib-more/xpassword"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xregexp"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
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

		goapidoc.NewRoutePath("GET", "/v1/auth/user", "当前登录用户").
			Tags("Authorization").
			Securities("Jwt").
			Responses(goapidoc.NewResponse(200, "_Result<UserDto>")),

		goapidoc.NewRoutePath("DELETE", "/v1/auth/logout", "注销").
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
	accountService *service.AccountService
	tokenService   *service.TokenService
	jwtService     *service.JwtService
}

func NewAuthController() *AuthController {
	return &AuthController{
		accountService: xdi.GetByNameForce(sn.SAccountService).(*service.AccountService),
		tokenService:   xdi.GetByNameForce(sn.STokenService).(*service.TokenService),
		jwtService:     xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
	}
}

// POST /v1/auth/register
func (a *AuthController) Register(c *gin.Context) *result.Result {
	pa := &param.RegisterParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	encrypted, err := xpassword.EncryptWithDefaultCost([]byte(pa.Password))
	if err != nil {
		return result.Error(exception.RegisterError).SetError(err, c)
	}

	status, err := a.accountService.Insert(pa.Email, string(encrypted))
	if status == xstatus.DbExisted {
		return result.Error(exception.EmailRegisteredError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.RegisterError).SetError(err, c)
	}

	return result.Created()
}

// POST /v1/auth/login
func (a *AuthController) Login(c *gin.Context) *result.Result {
	pa := &param.LoginParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	// get account by email / uid / username
	var account *po.Account
	var err error

	if xregexp.EmailRegex.MatchString(pa.Parameter) {
		account, err = a.accountService.QueryByEmail(pa.Parameter)
	} else if f := pa.Parameter[0]; f >= '0' && f <= '9' {
		uid, err := xnumber.Atou64(pa.Parameter)
		if err != nil {
			return result.Error(exception.RequestParamError).SetError(err, c)
		}
		account, err = a.accountService.QueryByUid(uid)
	} else {
		account, err = a.accountService.QueryByUsername(pa.Parameter)
	}

	if err != nil {
		return result.Error(exception.LoginError).SetError(err, c)
	} else if account == nil {
		return result.Error(exception.LoginParameterError)
	}

	// check password
	ok, err := xpassword.Check([]byte(pa.Password), []byte(account.Password))
	if err != nil {
		return result.Error(exception.LoginError).SetError(err, c)
	} else if !ok {
		return result.Error(exception.LoginParameterError)
	}

	// handle token
	token, err := a.jwtService.GenerateToken(account.Uid)
	if err != nil {
		return result.Error(exception.LoginError).SetError(err, c)
	}
	err = a.tokenService.Insert(token, account.Uid)
	if err != nil {
		return result.Error(exception.LoginError).SetError(err, c)
	}

	// reply
	res := dto.BuildLoginDto(account.User, token)
	return result.Ok().SetData(res)
}

// GET /v1/auth/user
func (a *AuthController) CurrentUser(c *gin.Context) *result.Result {
	user := a.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	res := dto.BuildUserDto(user)
	return result.Ok().SetData(res)
}

// DELETE /v1/auth/logout
func (a *AuthController) Logout(c *gin.Context) *result.Result {
	user := a.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	token := a.jwtService.GetToken(c)
	err := a.tokenService.Delete(token)
	if err != nil {
		return result.Error(exception.LogoutError).SetError(err, c)
	}

	return result.Ok()
}

// PUT /v1/auth/password
func (a *AuthController) UpdatePassword(c *gin.Context) *result.Result {
	user := a.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	pa := &param.UpdatePasswordParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	// get account
	account, err := a.accountService.QueryByUser(user)
	if err != nil {
		return result.Error(exception.WrongPasswordError).SetError(err, c)
	} else if account == nil {
		return result.Error(exception.UnAuthorizedError)
	}

	// check password
	ok, err := xpassword.Check([]byte(pa.Old), []byte(account.Password))
	if err != nil {
		return result.Error(exception.UpdatePasswordError).SetError(err, c)
	} else if !ok {
		return result.Error(exception.WrongPasswordError)
	}
	encrypted, err := xpassword.EncryptWithDefaultCost([]byte(pa.New))
	if err != nil {
		return result.Error(exception.UpdatePasswordError).SetError(err, c)
	}

	// update mysql and redis
	status, err := a.accountService.UpdatePassword(account.Uid, string(encrypted))
	if status == xstatus.DbNotFound {
		return result.Error(exception.UnAuthorizedError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UpdatePasswordError).SetError(err, c)
	}
	_ = a.tokenService.DeleteAll(account.Uid)

	// reply
	return result.Ok()
}
