package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
	"github.com/vidorg/vid_backend/src/util"
)

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

// @Router              /v1/auth/login [POST]
// @Summary             登录
// @Tag                 Authorization
// @Param               param body #LoginParam true "请求参数"
// @ResponseModel 200   #Result<LoginDto>
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

// @Router              /v1/auth/register [POST]
// @Summary             注册
// @Tag                 Authorization
// @Param               param body #RegisterParam true "请求参数"
// @ResponseModel 201   #Result<UserDto>
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
	if status == database.DbExisted {
		result.Error(exception.UsernameUsedError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.RegisterError).JSON(c)
		return
	}

	ret := dto.BuildUserDto(account.User)
	result.Created().SetData(ret).JSON(c)
}

// @Router              /v1/auth [GET]
// @Summary             当前登录用户
// @Tag                 Authorization
// @Security            Jwt
// @ResponseModel 200   #Result<UserDto>
func (a *AuthController) CurrentUser(c *gin.Context) {
	user := a.jwtService.GetContextUser(c)

	ret := dto.BuildUserDto(user)
	result.Ok().SetData(ret).JSON(c)
}

// @Router              /v1/auth/logout [POST]
// @Summary             注销
// @Tag                 Authorization
// @Security            Jwt
// @ResponseModel 200   #Result
func (a *AuthController) Logout(c *gin.Context) {
	token := a.jwtService.GetToken(c)
	ok := a.tokenService.Delete(token)
	if !ok {
		result.Error(exception.LogoutError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// @Router              /v1/auth/password [PUT]
// @Summary             修改密码
// @Tag                 Authorization
// @Security            Jwt
// @Param               param body #PassParam true "请求参数"
// @ResponseModel 200   #Result
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
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.UpdatePassError).JSON(c)
		return
	}
	_ = a.tokenService.DeleteAll(authUser.Uid)

	result.Ok().JSON(c)
}
