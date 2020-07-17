package controller

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/ahlib/xslice"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/service"
)

type UserController struct {
	Config        *config.Config            `di:"~"`
	Logger        *logrus.Logger            `di:"~"`
	Mappers       *xentity.EntityMappers    `di:"~"`
	JwtService    *service.JwtService       `di:"~"`
	UserService   *service.UserService      `di:"~"`
	SubService    *service.SubscribeService `di:"~"`
	VideoService  *service.VideoService     `di:"~"`
	SearchService *service.SearchService    `di:"~"`
}

func NewUserController(dic *xdi.DiContainer) *UserController {
	ctrl := &UserController{}
	dic.MustInject(ctrl)
	return ctrl
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
	pageOrder := param.BindPageOrder(c, u.Config)
	users, count := u.UserService.QueryAll(pageOrder)

	retDto := xcondition.First(u.Mappers.MapSlice(xslice.Sti(users), &dto.UserDto{}, dto.UserDtoShowAllOption())).([]*dto.UserDto)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, retDto).JSON(c)
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

	user := u.UserService.QueryByUid(uid)
	if user == nil {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}
	subscribingCnt, subscriberCnt, _ := u.SubService.QueryCountByUid(user.Uid)
	videoCnt, _ := u.VideoService.QueryCountByUid(user.Uid)
	extraInfo := &dto.UserExtraDto{
		SubscribingCount: subscribingCnt,
		SubscriberCount:  subscriberCnt,
		VideoCount:       videoCnt,
	}

	authUser := u.JwtService.GetContextUser(c)
	retDto := xcondition.First(u.Mappers.Map(user, &dto.UserDto{}, dto.UserDtoCheckUserOption(authUser))).(*dto.UserDto)
	result.Ok().
		PutData("user", retDto).
		PutData("extra", extraInfo).JSON(c)
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
		user := &po.User{}
		if !isSpec {
			user = u.JwtService.GetContextUser(c)
		} else {
			uid, ok := param.BindRouteId(c, "uid")
			if !ok {
				result.Error(exception.RequestParamError).JSON(c)
				return
			}
			user = u.UserService.QueryByUid(uid)
			if user == nil {
				result.Error(exception.UserNotFoundError).JSON(c)
				return
			}
		}
		// Update
		userParam := &param.UserParam{}
		if err := c.ShouldBind(userParam); err != nil {
			result.Error(exception.WrapValidationError(err)).PutData("error", err.Error()).JSON(c)
			return
		}

		_ = u.Mappers.MapProp(userParam, user)
		status := u.UserService.Update(user)
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

		retDto := xcondition.First(u.Mappers.Map(user, &dto.UserDto{}, dto.UserDtoShowAllOption())).(*dto.UserDto)
		result.Ok().SetData(retDto).JSON(c)
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
		var uid int32
		if !isSpec {
			uid = u.JwtService.GetContextUser(c).Uid
		} else {
			var ok bool
			uid, ok = param.BindRouteId(c, "uid")
			if !ok {
				result.Error(exception.RequestParamError).JSON(c)
				return
			}
		}
		// Delete
		status := u.UserService.Delete(uid)
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
