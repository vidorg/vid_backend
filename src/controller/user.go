package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/goapidoc"
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
)

func init() {
	goapidoc.AddPaths(
		goapidoc.NewPath("GET", "/v1/user", "查询所有用户").
			WithTags("User").
			WithSecurities("Jwt").
			WithParams(param.ADPage, param.ADLimit, param.ADOrder).
			WithResponses(goapidoc.NewResponse(200).WithType("_Result<_Page<UserDto>>")),

		goapidoc.NewPath("GET", "/v1/user/{uid}", "查询用户").
			WithTags("User").
			WithParams(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
			WithResponses(goapidoc.NewResponse(200).WithType("_Result<UserDetailDto>")),

		goapidoc.NewPath("PUT", "/v1/user/", "更新用户").
			WithTags("User").
			WithSecurities("Jwt").
			WithParams(goapidoc.NewBodyParam("param", "UserParam", true, "用户请求参数")).
			WithResponses(goapidoc.NewResponse(200).WithType("_Result<UserDto>")),

		goapidoc.NewPath("PUT", "/v1/user/admin/{uid}", "管理员更新用户").
			WithTags("User").
			// WithTags("Administration").
			WithSecurities("Jwt").
			WithParams(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				goapidoc.NewBodyParam("param", "UserParam", true, "用户请求参数"),
			).
			WithResponses(goapidoc.NewResponse(200).WithType("_Result<UserDto>")),

		goapidoc.NewPath("DELETE", "/v1/user/", "删除用户").
			WithTags("User").
			WithSecurities("Jwt").
			WithResponses(goapidoc.NewResponse(200).WithType("Result")),

		goapidoc.NewPath("PUT", "/v1/user/admin/{uid}", "管理员删除用户").
			WithTags("User").
			// WithTags("Administration").
			WithSecurities("Jwt").
			WithParams(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
			WithResponses(goapidoc.NewResponse(200).WithType("Result")),
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
func (u *UserController) QueryAllUsers(c *gin.Context) {
	pageOrder := param.BindPageOrder(c, u.config)
	users, total := u.userService.QueryAll(pageOrder)

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(pageOrder.Page, pageOrder.Limit, total, ret).JSON(c)
}

// GET /v1/user/:uid
func (u *UserController) QueryUser(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	user := u.userService.QueryByUid(uid)
	if user == nil {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	subscribingCnt, subscriberCnt, _ := u.subscribeService.QueryCountByUid(user.Uid)
	videoCnt, _ := u.videoService.QueryCountByUid(user.Uid)

	extra := dto.BuildUserExtraDto(subscribingCnt, subscriberCnt, videoCnt)
	ret := dto.BuildUserDetailDto(user, extra)
	result.Ok().SetData(ret).JSON(c)
}

// PUT /v1/user
// PUT /v1/user/admin/:uid
func (u *UserController) UpdateUser(isSpec bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get where parameter
		user := &po.User{}
		if !isSpec {
			user = u.jwtService.GetContextUser(c)
		} else {
			uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				result.Error(exception.RequestParamError).JSON(c)
				return
			}
			user = u.userService.QueryByUid(uid)
			if user == nil {
				result.Error(exception.UserNotFoundError).JSON(c)
				return
			}
		}

		// Update
		userParam := &param.UserParam{}
		if err := c.ShouldBind(userParam); err != nil {
			result.Error(exception.WrapValidationError(err)).JSON(c)
			return
		}

		param.MapUserParam(userParam, user)
		status := u.userService.Update(user)
		if status == database.DbNotFound {
			result.Error(exception.UserNotFoundError).JSON(c)
			return
		} else if status == database.DbExisted {
			result.Error(exception.UsernameUsedError).JSON(c)
			return
		} else if status == database.DbFailed {
			result.Error(exception.UserUpdateError).JSON(c)
			return
		}

		ret := dto.BuildUserDto(user)
		result.Ok().SetData(ret).JSON(c)
	}
}

// DELETE /v1/user
// DELETE /v1/user/admin/:uid
func (u *UserController) DeleteUser(isSpec bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get delete uid param
		var uid int32
		if !isSpec {
			uid = u.jwtService.GetContextUser(c).Uid
		} else {
			var ok bool
			uid, ok = param.BindRouteId(c, "uid")
			if !ok {
				result.Error(exception.RequestParamError).JSON(c)
				return
			}
		}

		// Delete
		status := u.userService.Delete(uid)
		if status == database.DbNotFound {
			result.Error(exception.UserNotFoundError).JSON(c)
			return
		} else if status == database.DbFailed {
			result.Error(exception.UserDeleteError).JSON(c)
			return
		}

		result.Ok().JSON(c)
	}
}
