package controller

import (
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
	goapidoc.AddRoutePaths(
		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/favorite", "query user favorites").
			Tags("Favorite").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedIsBlock, _adNeedVideoCount, _adNeedAuthor,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/video/{vid}/favored", "query video favored users").
			Tags("Favorite").
			Params(
				goapidoc.NewPathParam("vid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedIsBlock, _adNeedVideoCount,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),
	)
}

type FavoriteController struct {
	config          *config.Config
	jwtService      *service.JwtService
	favoriteService *service.FavoriteService
	common          *CommonController
}

func NewFavoriteController() *FavoriteController {
	return &FavoriteController{
		config:          xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:      xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		favoriteService: xdi.GetByNameForce(sn.SFavoriteService).(*service.FavoriteService),
		common:          xdi.GetByNameForce(sn.SCommonController).(*CommonController),
	}
}

// /v1/user/:uid/favorite
func (f *FavoriteController) QueryFavorites(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, f.config)

	videos, total, err := f.favoriteService.QueryFavorites(uid, pp)
	if err != nil {
		return result.Error(exception.GetBlockingListError).SetError(err, c)
	} else if videos == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := f.jwtService.GetContextUser(c)
	authors, extras, err := f.common.getVideosAuthor(c, authUser, videos)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
	}

	res := dto.BuildVideoDtos(videos)
	for idx, video := range res {
		video.Author = dto.BuildUserDto(authors[idx])
		if video.Author != nil {
			video.Author.Extra = extras[idx]
		}
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// /v1/video/:vid/favored
func (f *FavoriteController) QueryFavoreds(c *gin.Context) *result.Result {
	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, f.config)

	users, total, err := f.favoriteService.QueryFavoreds(vid, pp)
	if err != nil {
		return result.Error(exception.GetBlockingListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.VideoNotFoundError)
	}

	authUser := f.jwtService.GetContextUser(c)
	extras, err := f.common.getUsersExtra(c, authUser, users)
	if err != nil {
		return result.Error(exception.QueryUserError).SetError(err, c)
	}

	res := dto.BuildUserDtos(users)
	for idx, user := range res {
		user.Extra = extras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}
