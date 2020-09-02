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
		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/subscriber", "query user subscribers").
			Tags("Subscribe").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedBlockCount, _adNeedVideoCount,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/subscribing", "query user subscribings").
			Tags("Subscribe").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedBlockCount, _adNeedVideoCount,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("POST", "/v1/user/subscribing/{uid}", "subscribe user").
			Tags("Subscribe").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int64", true, "user id")).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/user/subscribing/{uid}", "unsubscribe user").
			Tags("Subscribe").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int64", true, "user id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type SubscribeController struct {
	config           *config.Config
	jwtService       *service.JwtService
	userService      *service.UserService
	subscribeService *service.SubscribeService
	common           *CommonController
}

func NewSubscribeController() *SubscribeController {
	return &SubscribeController{
		config:           xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:       xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		userService:      xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
		common:           xdi.GetByNameForce(sn.SCommonController).(*CommonController),
	}
}

// GET /v1/user/:uid/subscriber
func (s *SubscribeController) QuerySubscribers(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, s.config)

	users, total, err := s.subscribeService.QuerySubscribers(uid, pp)
	if err != nil {
		return result.Error(exception.GetSubscriberListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := s.jwtService.GetContextUser(c)
	extras, err := s.common.getUsersExtra(c, authUser, users)
	if err != nil {
		return result.Error(exception.QueryUserError).SetError(err, c)
	}

	res := dto.BuildUserDtos(users)
	for idx, user := range res {
		user.Extra = extras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/user/{uid}/subscribing
func (s *SubscribeController) QuerySubscribings(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, s.config)

	users, total, err := s.subscribeService.QuerySubscribings(uid, pp)
	if err != nil {
		return result.Error(exception.GetSubscribingListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := s.jwtService.GetContextUser(c)
	extras, err := s.common.getUsersExtra(c, authUser, users)
	if err != nil {
		return result.Error(exception.QueryUserError).SetError(err, c)
	}

	res := dto.BuildUserDtos(users)
	for idx, user := range res {
		user.Extra = extras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// POST /v1/user/subscribing/:uid
func (s *SubscribeController) SubscribeUser(c *gin.Context) *result.Result {
	user := s.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	if user.Uid == uid {
		return result.Error(exception.SubscribeSelfError)
	}

	status, err := s.subscribeService.InsertSubscribe(user.Uid, uid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbExisted {
		return result.Error(exception.AlreadySubscribingError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.SubscribeError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/user/subscribing/:uid
func (s *SubscribeController) UnSubscribeUser(c *gin.Context) *result.Result {
	user := s.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	status, err := s.subscribeService.DeleteSubscribe(user.Uid, uid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbTagA {
		return result.Error(exception.NotSubscribeYetError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UnSubscribeError).SetError(err, c)
	}

	return result.Ok()
}
