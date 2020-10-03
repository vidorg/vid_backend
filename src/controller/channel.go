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
		goapidoc.NewRoutePath("GET", "/v1/channel", "query all channels").
			Tags("Channel", "Administration").
			Securities("Jwt").
			Params(
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ChannelDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/channel", "query channels from user").
			Tags("Channel").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<ChannelDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/channel/{cid}", "query a channel").
			Tags("Channel").
			Params(
				goapidoc.NewPathParam("cid", "integer#int64", true, "channel id"),
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<ChannelDto>")),

		goapidoc.NewRoutePath("POST", "/v1/channel", "create a channel").
			Tags("Channel").
			Securities("Jwt").
			Params(goapidoc.NewBodyParam("param", "InsertChannelParam", true, "create channel parameter")).
			Responses(goapidoc.NewResponse(201, "Result")),

		goapidoc.NewRoutePath("PUT", "/v1/channel/{cid}", "update a channel").
			Tags("Channel").
			Securities("Jwt").
			Params(
				goapidoc.NewPathParam("cid", "integer#int64", true, "channel id"),
				goapidoc.NewBodyParam("param", "UpdateChannelParam", true, "update channel parameter"),
			).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/channel/{cid}", "delete a channel").
			Tags("Channel").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("cid", "integer#int64", true, "channel id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type ChannelController struct {
	config         *config.Config
	common         *CommonController
	jwtService     *service.JwtService
	channelService *service.ChannelService
}

func NewChannelController() *ChannelController {
	return &ChannelController{
		config:         xdi.GetByNameForce(sn.SConfig).(*config.Config),
		common:         xdi.GetByNameForce(sn.SCommonController).(*CommonController),
		jwtService:     xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		channelService: xdi.GetByNameForce(sn.SChannelService).(*service.ChannelService),
	}
}

// GET /v1/channel
func (ch *ChannelController) QueryAllChannels(c *gin.Context) *result.Result {
	pp := param.BindPageOrder(c, ch.config)
	channels, total, err := ch.channelService.QueryAll(pp)
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	}

	authUser := ch.jwtService.GetContextUser(c)
	authors, userExtras, err := ch.common.getChannelAuthors(c, authUser, channels)
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	}
	channelExtras, err := ch.common.getChannelExtras(c, authUser, channels)
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	}

	res := dto.BuildChannelDtos(channels)
	for idx, channel := range res {
		channel.Author = dto.BuildUserDto(authors[idx])
		if channel.Author != nil {
			channel.Author.Extra = userExtras[idx]
		}
		channel.Extra = channelExtras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/user/:uid/channel
func (ch *ChannelController) QueryChannelsByUid(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, ch.config)

	channels, total, err := ch.channelService.QueryByUid(uid, pp)
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	} else if channels == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := ch.jwtService.GetContextUser(c)
	authors, userExtras, err := ch.common.getChannelAuthors(c, authUser, channels)
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	}
	channelExtras, err := ch.common.getChannelExtras(c, authUser, channels)
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	}

	res := dto.BuildChannelDtos(channels)
	for idx, channel := range res {
		channel.Author = dto.BuildUserDto(authors[idx])
		if channel.Author != nil {
			channel.Author.Extra = userExtras[idx]
		}
		channel.Extra = channelExtras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/channel/:cid
func (ch *ChannelController) QueryChannelByCid(c *gin.Context) *result.Result {
	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	channel, err := ch.channelService.QueryByCid(cid)
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	} else if channel == nil {
		return result.Error(exception.ChannelNotFoundError)
	}

	authUser := ch.jwtService.GetContextUser(c)
	authors, userExtras, err := ch.common.getChannelAuthors(c, authUser, []*po.Channel{channel})
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	}
	channelExtras, err := ch.common.getChannelExtras(c, authUser, []*po.Channel{channel})
	if err != nil {
		return result.Error(exception.QueryChannelError).SetError(err, c)
	}

	res := dto.BuildChannelDto(channel)
	res.Author = dto.BuildUserDto(authors[0])
	if res.Author != nil {
		res.Author.Extra = userExtras[0]
	}
	res.Extra = channelExtras[0]
	return result.Ok().SetData(res)
}

// POST /v1/channel
func (ch *ChannelController) InsertChannel(c *gin.Context) *result.Result {
	user := ch.jwtService.GetContextUser(c)
	pa := &param.InsertChannelParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	status, err := ch.channelService.Insert(pa, user.Uid)
	if status == xstatus.DbExisted {
		return result.Error(exception.ChannelNameUsedError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.InsertChannelError).SetError(err, c)
	}

	return result.Created()
}

// PUT /v1/channel/:cid
func (ch *ChannelController) UpdateChannel(c *gin.Context) *result.Result {
	user := ch.jwtService.GetContextUser(c)
	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pa := &param.UpdateChannelParam{}
	if err := c.ShouldBind(pa); err != nil {
		return result.Error(exception.WrapValidationError(err)).SetError(err, c)
	}

	channel, err := ch.channelService.QueryByCid(cid)
	if err != nil {
		return result.Error(exception.UpdateChannelError).SetError(err, c)
	} else if channel == nil {
		return result.Error(exception.ChannelNotFoundError)
	} else if channel.AuthorUid != user.Uid {
		return result.Error(exception.ChannelPermissionError)
	}

	status, err := ch.channelService.Update(cid, pa)
	if status == xstatus.DbNotFound {
		return result.Error(exception.ChannelNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UpdateChannelError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/channel/:cid
func (ch *ChannelController) DeleteChannel(c *gin.Context) *result.Result {
	user := ch.jwtService.GetContextUser(c)
	cid, err := param.BindRouteId(c, "cid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	channel, err := ch.channelService.QueryByCid(cid)
	if err != nil {
		return result.Error(exception.DeleteChannelError).SetError(err, c)
	} else if channel == nil {
		return result.Error(exception.ChannelNotFoundError)
	} else if channel.AuthorUid != user.Uid {
		return result.Error(exception.ChannelPermissionError)
	}

	status, err := ch.channelService.Delete(cid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.ChannelNotFoundError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.DeleteChannelError).SetError(err, c)
	}

	return result.Ok()
}
