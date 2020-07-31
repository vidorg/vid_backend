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
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

func init() {
	goapidoc.AddPaths(
		goapidoc.NewPath("GET", "/v1/user/{uid}/subscriber", "查询用户粉丝").
			WithTags("Subscribe").
			WithParams(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				param.ADPage, param.ADLimit, param.ADOrder,
			).
			WithResponses(goapidoc.NewResponse(200).WithType("_Result<_Page<UserDto>>")),

		goapidoc.NewPath("GET", "/v1/user/{uid}/subscribing", "查询用户关注").
			WithTags("Subscribe").
			WithParams(
				goapidoc.NewPathParam("uid", "integer#int32", true, "用户id"),
				param.ADPage, param.ADLimit, param.ADOrder,
			).
			WithResponses(goapidoc.NewResponse(200).WithType("_Result<_Page<UserDto>>")),

		goapidoc.NewPath("PUT", "/v1/user/subscribing/{uid}", "关注用户").
			WithTags("Subscribe").
			WithSecurities("Jwt").
			WithParams(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
			WithResponses(goapidoc.NewResponse(200).WithType("Result")),

		goapidoc.NewPath("DELETE", "/v1/user/subscribing/{uid}", "取消关注用户").
			WithTags("Subscribe").
			WithSecurities("Jwt").
			WithParams(goapidoc.NewPathParam("uid", "integer#int32", true, "用户id")).
			WithResponses(goapidoc.NewResponse(200).WithType("Result")),
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
	pageOrder := param.BindPageOrder(c, s.config)

	users, total, status := s.subscribeService.QuerySubscriberUsers(uid, pageOrder)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(pageOrder.Page, pageOrder.Limit, total, ret).JSON(c)
}

// GET /v1/user/{uid}/subscribing
func (s *SubscribeController) QuerySubscribingUsers(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pageOrder := param.BindPageOrder(c, s.config)

	users, total, status := s.subscribeService.QuerySubscribingUsers(uid, pageOrder)
	if status == xstatus.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(pageOrder.Page, pageOrder.Limit, total, ret).JSON(c)
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
