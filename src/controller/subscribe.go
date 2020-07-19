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
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

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

// @Router              /v1/user/{uid}/subscriber [GET]
// @Summary             查询用户粉丝
// @Tag                 Subscribe
// @Template            Order Page
// @Param               uid path integer true "用户id"
// @ResponseModel 200   #Result<Page<UserDto>>
func (s *SubscribeController) QuerySubscriberUsers(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pageOrder := param.BindPageOrder(c, s.config)

	users, count, status := s.subscribeService.QuerySubscriberUsers(uid, pageOrder)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, ret).JSON(c)
}

// @Router              /v1/user/{uid}/subscribing [GET]
// @Summary             查询用户关注
// @Tag                 Subscribe
// @Template            Order Page
// @Param               uid path integer true "用户id"
// @ResponseModel 200   #Result<Page<UserDto>>
func (s *SubscribeController) QuerySubscribingUsers(c *gin.Context) {
	uid, ok := param.BindRouteId(c, "uid")
	if !ok {
		result.Error(exception.RequestParamError).JSON(c)
		return
	}
	pageOrder := param.BindPageOrder(c, s.config)

	users, count, status := s.subscribeService.QuerySubscribingUsers(uid, pageOrder)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	}

	ret := dto.BuildUserDtos(users)
	result.Ok().SetPage(count, pageOrder.Page, pageOrder.Limit, ret).JSON(c)
}

// @Router              /v1/user/subscribing [PUT]
// @Summary             关注用户
// @Tag                 Subscribe
// @Security            Jwt
// @Param               param body #SubscribeParam true "请求参数"
// @ResponseModel 200   #Result
func (s *SubscribeController) SubscribeUser(c *gin.Context) {
	authUser := s.jwtService.GetContextUser(c)
	subParam := &param.SubscribeParam{}
	if err := c.ShouldBind(subParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}
	if authUser.Uid == subParam.To {
		result.Error(exception.SubscribeSelfError).JSON(c)
		return
	}

	status := s.subscribeService.SubscribeUser(authUser.Uid, subParam.To)
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
// @Summary             取消关注用户
// @Tag                 Subscribe
// @Security            Jwt
// @Param               param body #SubscribeParam true "请求参数"
// @ResponseModel 200   #Result
func (s *SubscribeController) UnSubscribeUser(c *gin.Context) {
	authUser := s.jwtService.GetContextUser(c)
	subParam := &param.SubscribeParam{}
	if err := c.ShouldBind(subParam); err != nil {
		result.Error(exception.WrapValidationError(err)).JSON(c)
		return
	}

	status := s.subscribeService.UnSubscribeUser(authUser.Uid, subParam.To)
	if status == database.DbNotFound {
		result.Error(exception.UserNotFoundError).JSON(c)
		return
	} else if status == database.DbFailed {
		result.Error(exception.UnSubscribeError).JSON(c)
		return
	}

	result.Ok().JSON(c)
}
