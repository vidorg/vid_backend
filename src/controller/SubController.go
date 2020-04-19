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
	"github.com/vidorg/vid_backend/src/service"
)

type SubController struct {
	Config      *config.ServerConfig      `di:"~"`
	Logger      *logrus.Logger            `di:"~"`
	Mappers     *xentity.EntityMappers    `di:"~"`
	JwtService  *service.JwtService       `di:"~"`
	UserService *service.UserService      `di:"~"`
	SubService  *service.SubscribeService `di:"~"`
}

func NewSubController(dic *xdi.DiContainer) *SubController {
	ctrl := &SubController{}
	dic.MustInject(ctrl)
	return ctrl
}

// @Router              /v1/user/{uid}/subscriber [GET]
// @Template            Page ParamA Order Page
// @Summary             查询用户粉丝
// @Tag                 Subscribe
// @Param               uid path integer true "用户id"
// @ResponseDesc 404    "user not found"
// @ResponseModel 200   #Result<Page<UserDto>>
func (s *SubController) QuerySubscriberUsers(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pageOrder := param.BindPageOrder(c, s.Config)

	users, count, status := s.SubService.QuerySubscriberUsers(uid, pageOrder)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	retDto := xcondition.First(s.Mappers.MapSlice(xslice.Sti(users), &dto.UserDto{})).([]*dto.UserDto)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, retDto).JSON(c)
}

// @Router              /v1/user/{uid}/subscribing [GET]
// @Template            Page ParamA Order Page
// @Summary             查询用户关注
// @Tag                 Subscribe
// @Param               uid path integer true "用户id"
// @ResponseDesc 404    "user not found"
// @ResponseModel 200   #Result<Page<UserDto>>
func (s *SubController) QuerySubscribingUsers(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pageOrder := param.BindPageOrder(c, s.Config)

	users, count, status := s.SubService.QuerySubscribingUsers(uid, pageOrder)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	retDto := xcondition.First(s.Mappers.MapSlice(xslice.Sti(users), &dto.UserDto{})).([]*dto.UserDto)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, retDto).JSON(c)
}

// @Router              /v1/user/subscribing [PUT]
// @Security            Jwt
// @Template            Auth Param
// @Summary             关注用户
// @Tag                 Subscribe
// @Param               param body #SubParam true "请求参数"
// @ResponseDesc 400    "subscribe oneself invalid"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "subscribe failed"
// @ResponseModel 200   #Result
func (s *SubController) SubscribeUser(c *gin.Context) {
	authUser := s.JwtService.GetContextUser(c)
	subParam := &param.SubParam{}
	if err := c.ShouldBind(subParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}
	if authUser.Uid == subParam.Uid {
		result.Error(exception.SubscribeSelfError).JSON(c)
		return
	}

	status := s.SubService.SubscribeUser(authUser.Uid, subParam.Uid)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.SubscribeError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// @Router              /v1/user/subscribing [DELETE]
// @Security            Jwt
// @Template            Auth Param
// @Summary             取消关注用户
// @Tag                 Subscribe
// @Param               param body #SubParam true "请求参数"
// @ResponseDesc 404    "user not found"
// @ResponseDesc 500    "unsubscribe failed"
// @ResponseModel 200   #Result
func (s *SubController) UnSubscribeUser(c *gin.Context) {
	authUser := s.JwtService.GetContextUser(c)
	subParam := &param.SubParam{}
	if err := c.ShouldBind(subParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	status := s.SubService.UnSubscribeUser(authUser.Uid, subParam.Uid)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.UnSubscribeError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}
