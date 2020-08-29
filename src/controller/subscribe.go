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
		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/subscriber", "查询用户粉丝").
			Tags("Subscribe").
			Params(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				param.ADPage, param.ADLimit, param.ADOrder,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/subscribing", "查询用户关注").
			Tags("Subscribe").
			Params(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				param.ADPage, param.ADLimit, param.ADOrder,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("PUT", "/v1/user/subscribing/{uid}", "关注用户").
			Tags("Subscribe").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/user/subscribing/{uid}", "取消关注用户").
			Tags("Subscribe").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type SubscribeController struct {
	config           *config.Config
	jwtService       *service.JwtService
	userService      *service.UserService
	subscribeService *service.SubscribeService
}

func NewSubscribeController() *SubscribeController {
	return &SubscribeController{
		config:           xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:       xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		userService:      xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
	}
}

// GET /v1/user/{uid}/subscriber
func (s *SubscribeController) QuerySubscriberUsers(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pp := param.BindPageOrder(c, s.config)

	users, total, status := s.subscribeService.QuerySubscriberUsers(uid, pp)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(pp.Page, pp.Limit, total, ret).JSON(c)
}

// GET /v1/user/{uid}/subscribing
func (s *SubscribeController) QuerySubscribingUsers(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pp := param.BindPageOrder(c, s.config)

	users, total, status := s.subscribeService.QuerySubscribingUsers(uid, pp)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(pp.Page, pp.Limit, total, ret).JSON(c)
}

// PUT /v1/user/subscribing/:uid
func (s *SubscribeController) SubscribeUser(c *gin.Context) {
	authUser := s.jwtService.GetContextUser(c)
	to, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	if authUser.Uid == to {
		result.Error(exception.SubscribeSelfError).JSON(c)
		return
	}

	status := s.subscribeService.SubscribeUser(authUser.Uid, to)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == xstatus.DbFailed {
		result.Error(exception.SubscribeError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}

// PUT /v1/user/subscribing/:uid
func (s *SubscribeController) UnSubscribeUser(c *gin.Context) {
	authUser := s.jwtService.GetContextUser(c)
	to, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}

	status := s.subscribeService.UnSubscribeUser(authUser.Uid, to)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == xstatus.DbFailed {
		result.Error(exception.UnSubscribeError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}
