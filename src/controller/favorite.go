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
		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/favorite", "query user favorites").
			Tags("Favorite").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedVideoChannel, _adNeedFavoredCount, _adNeedIsFavorite,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<VideoDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/video/{vid}/favored", "query video favored users").
			Tags("Favorite").
			Params(
				goapidoc.NewPathParam("vid", "integer#int64", true, "vid id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedVideoChannel, _adNeedFavoredCount, _adNeedIsFavorite,
				_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("POST", "/v1/user/favorite/{vid}", "add video to favorite").
			Tags("Favorite").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("vid", "integer#int64", true, "vid id")).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/user/favorite/{vid}", "remove video from favorite").
			Tags("Favorite").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("vid", "integer#int64", true, "vid id")).
			Responses(goapidoc.NewResponse(200, "Result")),
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
		return result.Error(exception.GetFavoriteListError).SetError(err, c)
	} else if videos == nil {
		return result.Error(exception.UserNotFoundError)
	}

	res := dto.BuildVideoDtos(videos)
	err = f.common.PreLoadVideos(c, f.jwtService.GetContextUser(c), videos, res)
	if err != nil {
		return result.Error(exception.QueryVideoError).SetError(err, c)
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
		return result.Error(exception.GetFavoredListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.VideoNotFoundError)
	}

	res := dto.BuildUserDtos(users)
	err = f.common.PreLoadUsers(c, f.jwtService.GetContextUser(c), users, res)
	if err != nil {
		return result.Error(exception.GetFavoredListError).SetError(err, c)
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// POST /v1/user/favorite/:vid
func (f *FavoriteController) AddFavorite(c *gin.Context) *result.Result {
	user := f.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	status, err := f.favoriteService.InsertFavorite(user.Uid, vid)
	if status == xstatus.DbTagB {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbTagC {
		return result.Error(exception.VideoNotFoundError)
	} else if status == xstatus.DbExisted {
		return result.Error(exception.AlreadyInFavoriteError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.AddToFavoriteError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/user/favorite/:vid
func (f *FavoriteController) RemoveFavorite(c *gin.Context) *result.Result {
	user := f.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	vid, err := param.BindRouteId(c, "vid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	status, err := f.favoriteService.DeleteFavorite(user.Uid, vid)
	if status == xstatus.DbTagB {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbTagC {
		return result.Error(exception.VideoNotFoundError)
	} else if status == xstatus.DbTagA {
		return result.Error(exception.NotInFavoriteYetError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.RemoveFromFavoriteError).SetError(err, c)
	}

	return result.Ok()
}
