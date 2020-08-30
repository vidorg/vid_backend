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

		goapidoc.NewRoutePath("PUT", "/v1/user/admin/{uid}", "管理员更新用户").
			Tags("User", "Administration").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				goapidoc.NewBodyParam("param", "UpdateUserParam", true, "用户请求参数"),
			).
			Responses(goapidoc.NewResponse(200, "_Result<UserDto>")),

		goapidoc.NewRoutePath("DELETE", "/v1/user/", "删除用户").
			Tags("User").
			Securities("Jwt").
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("PUT", "/v1/user/admin/{uid}", "管理员删除用户").
			Tags("User", "Administration").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
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
func (u *UserController) QueryAllUsers(c *gin.Context) *result.Result {
	pp := param.BindPageOrder(c, u.config)
	users, total := u.userService.QueryAll(pp)

	ret := dto.BuildUserDtos(users)
	return result.Ok().SetPage(pp.Page, pp.Limit, total, ret)
}

// GET /v1/user/:uid
func (u *UserController) QueryUser(c *gin.Context) *result.Result {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		return result.Error(exception.RequestParamError)
	}

	user := u.userService.QueryByUid(uid)
	if user == nil {
		return result.Error(exception.UserNotFoundError)
	}

	subscribingCnt, subscriberCnt, _ := u.subscribeService.QueryCountByUid(user.Uid)
	videoCnt, _ := u.videoService.QueryCountByUid(user.Uid)

	extra := dto.BuildUserExtraDto(subscribingCnt, subscriberCnt, videoCnt)
	ret := dto.BuildUserDetailDto(user, extra)
	return result.Ok().SetData(ret)
}

// PUT /v1/user
// PUT /v1/user/admin/:uid
func (u *UserController) UpdateUser(isSpec bool) func(c *gin.Context) *result.Result {
	return func(c *gin.Context) *result.Result {
		// get where parameter
		user := &po.User{}
		if !isSpec {
			user = u.jwtService.GetContextUser(c)
		} else {
			uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				return result.Error(exception.RequestParamError)
			}
			user = u.userService.QueryByUid(uid)
			if user == nil {
				return result.Error(exception.UserNotFoundError)
			}
		}

		// Update
		userParam := &param.UpdateUserParam{}
		if err := c.ShouldBind(userParam); err != nil {
			return result.Error(exception.WrapValidationError(err))
		}

		param.MapUserParam(userParam, user)
		status := u.userService.Update(user)
		if status == xstatus.DbNotFound {
			return result.Error(exception.UserNotFoundError)
		} else if status == xstatus.DbExisted {
			return result.Error(exception.UsernameUsedError)
		} else if status == xstatus.DbFailed {
			return result.Error(exception.UserUpdateError)
		}

		ret := dto.BuildUserDto(user)
		return result.Ok().SetData(ret)
	}
}

// DELETE /v1/user
// DELETE /v1/user/admin/:uid
func (u *UserController) DeleteUser(isSpec bool) func(c *gin.Context) *result.Result {
	return func(c *gin.Context) *result.Result {
		// get delete uid param
		var uid int32
		if !isSpec {
			uid = u.jwtService.GetContextUser(c).Uid
		} else {
			var ok bool
			uid, ok = param.BindRouteId(c, "uid")
			if !ok {
				return result.Error(exception.RequestParamError)
			}
		}

		// Delete
		status := u.userService.Delete(uid)
		if status == xstatus.DbNotFound {
			return result.Error(exception.UserNotFoundError)
		} else if status == xstatus.DbFailed {
			return result.Error(exception.UserDeleteError)
		}

		return result.Ok()
	}
}
