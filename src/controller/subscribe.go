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
		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/subscribing", "query user subscribings").
			Tags("Subscribe").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ChannelDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/channel/{cid}/subscriber", "query channel subscriber users").
			Tags("Subscribe").
			Params(
				goapidoc.NewPathParam("cid", "integer#int64", true, "channel id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("POST", "/v1/user/subscribe/{cid}", "subscribe channel").
			Tags("Subscribe").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("cid", "integer#int64", true, "channel id")).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/user/subscribe/{cid}", "unsubscribe channel").
			Tags("Subscribe").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("cid", "integer#int64", true, "channel id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type SubscribeController struct {
	config           *config.Config
	jwtService       *service.JwtService
	subscribeService *service.SubscribeService
	common           *CommonController
}

func NewSubscribeController() *SubscribeController {
	return &SubscribeController{
		config:           xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:       xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
		common:           xdi.GetByNameForce(sn.SCommonController).(*CommonController),
	}
}

// GET /v1/user/:uid/subscribing
func (s *SubscribeController) QuerySubscribings(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, s.config)

	channels, total, err := s.subscribeService.QuerySubscribings(uid, pp)
	if err != nil {
		return result.Error(exception.GetSubscribingListError).SetError(err, c)
	} else if channels == nil {
		return result.Error(exception.UserNotFoundError)
	}

	res := dto.BuildChannelDtos(channels)
	err = s.common.PreLoadChannels(c, s.jwtService.GetContextUser(c), channels, res)
	if err != nil {
		return result.Error(exception.GetSubscribingListError).SetError(err, c)
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/channel/:cid/subscriber
func (s *SubscribeController) QuerySubscribers(c *gin.Context) *result.Result {
	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, s.config)

	users, total, err := s.subscribeService.QuerySubscribers(cid, pp)
	if err != nil {
		return result.Error(exception.GetSubscriberListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.ChannelNotFoundError)
	}

	res := dto.BuildUserDtos(users)
	err = s.common.PreLoadUsers(c, s.jwtService.GetContextUser(c), users, res)
	if err != nil {
		return result.Error(exception.GetSubscriberListError).SetError(err, c)
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// POST /v1/user/subscribe/:cid
func (s *SubscribeController) SubscribeChannel(c *gin.Context) *result.Result {
	user := s.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	status, err := s.subscribeService.SubscribeChannel(user.Uid, cid)
	if status == xstatus.DbTagB {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbTagC {
		return result.Error(exception.ChannelNotFoundError)
	} else if status == xstatus.DbExisted {
		return result.Error(exception.AlreadySubscribedError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.SubscribeChannelError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/user/subscribe/:cid
func (s *SubscribeController) UnsubscribeChannel(c *gin.Context) *result.Result {
	user := s.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	status, err := s.subscribeService.UnsubscribeChannel(user.Uid, cid)
	if status == xstatus.DbTagB {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbTagC {
		return result.Error(exception.ChannelNotFoundError)
	} else if status == xstatus.DbTagA {
		return result.Error(exception.NotSubscribeYetError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UnsubscribeChannelError).SetError(err, c)
	}

	return result.Ok()
}
