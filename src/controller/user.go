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
)

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

// @Router              /v1/user [GET]
// @Summary             查询所有用户
// @Description         管理员权限，此处可见用户手机号码
// @Tag                 User
// @Tag                 Administration
// @Security            Jwt
// @Template            Order Page
// @ResponseModel 200   #Result<Page<UserDto>>
func (u *UserController) QueryAllUsers(c *gin.Context) {
	pageOrder := param.BindPageOrder(c, u.config)
	users, count := u.userService.QueryAll(pageOrder)

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, ret).JSON(c)
}

// @Router              /v1/user/{uid} [GET]
// @Summary             查询用户
// @Description         此处用户本人可见手机号码，管理员不受限制
// @Tag                 User
// @Param               uid path integer true "用户id"
// @ResponseModel 200   #Result<UserExtraDto>
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

// @Router              /v1/user [PUT]
// @Summary             更新用户
// @Tag                 User
// @Security            Jwt
// @Param               param body #UserParam true "请求参数"
// @ResponseModel 200   #Result<UserDto>
//
// @Router              /v1/user/admin/{uid} [PUT]
// @Summary             更新用户
// @Description         管理员权限
// @Tag                 User
// @Tag                 Administration
// @Security            Jwt
// @Param               uid   path integer    true "用户id"
// @Param               param body #UserParam true "请求参数"
// @ResponseModel 200   #Result<UserDto>
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

// @Router              /v1/user [DELETE]
// @Security            Jwt
// @Summary             删除用户
// @Tag                 User
// @ResponseModel 200   #Result
//
// @Router              /v1/user/admin/{uid} [DELETE]
// @Security            Jwt
// @Summary             删除用户
// @Description         管理员权限
// @Tag                 User
// @Tag                 Administration
// @Param               uid path integer true "用户id"
// @ResponseModel 200   #Result
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
