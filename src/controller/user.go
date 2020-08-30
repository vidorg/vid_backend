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
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func init() {
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("GET", "/v1/user", "查询所有用户").
			Tags("User").
			Securities("Jwt").
			Params(param.ADPage, param.ADLimit, param.ADOrder).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/user/{uid}", "查询用户").
			Tags("User").
			Params(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
			Responses(goapidoc.NewResponse(200, "_Result<UserDetailDto>")),

		goapidoc.NewRoutePath("PUT", "/v1/user/", "更新用户").
			Tags("User").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("param", "UpdateUserParam", true, "用户请求参数")).
			Responses(goapidoc.NewResponse(200, "_Result<UserDto>")),

		goapidoc.NewRoutePath("DELETE", "/v1/user/", "删除用户").
			Tags("User").
			Securities("Jwt").
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type UserController struct {
	config           *config.Config
	jwtService       *service.JwtService
	userService      *service.UserService
	subscribeService *service.SubscribeService
	videoService     *service.VideoService
}

func NewUserController() *UserController {
	return &UserController{
		config:           xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:       xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		userService:      xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
		videoService:     xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
	}
}

// GET /v1/user
func (u *UserController) QueryAll(c *gin.Context) *result.Result {
	user := u.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	pp := param.BindPageOrder(c, u.config)
	users, total, err := u.userService.QueryAll(pp)
	if err != nil {
		return result.Error(exception.QueryUserError).SetError(err, c)
	}

	res := dto.BuildUserDtos(users)
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/user/uid/:uid
func (u *UserController) QueryByUid(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	user, err := u.userService.QueryByUid(uid)
	if err != nil {
		result.Error(exception.QueryUserError).SetError(err, c)
	} else if user == nil {
		return result.Error(exception.UserNotFoundError)
	}

	// TODO extra

	res := dto.BuildUserDto(user)
	return result.Ok().SetData(res)
}

// GET /v1/user/username/:uid
func (u *UserController) QueryByUsername(c *gin.Context) *result.Result {
	username := c.Param("username")
	if username == "" {
		return result.Error(exception.RequestParamError)
	}

	user, err := u.userService.QueryByUsername(username)
	if err != nil {
		result.Error(exception.QueryUserError).SetError(err, c)
	} else if user == nil {
		return result.Error(exception.UserNotFoundError)
	}

	// TODO extra

	res := dto.BuildUserDto(user)
	return result.Ok().SetData(res)
}

// PUT /v1/user
func (u *UserController) Update(c *gin.Context) *result.Result {
	user := u.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	pa := &param.UpdateUserParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	status, err := u.userService.Update(user.Uid, pa)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbExisted {
		return result.Error(exception.UsernameUsedError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UpdateUserError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/user
func (u *UserController) Delete(c *gin.Context) *result.Result {
	user := u.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	status, err := u.userService.Delete(user.Uid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.DeleteUserError).SetError(err, c)
	}

	return result.Ok()
}
