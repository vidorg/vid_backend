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
		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/follower", "query user followers").
			Tags("Follow").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("GET", "/v1/user/{uid}/following", "query user followings").
			Tags("Follow").
			Params(
				goapidoc.NewPathParam("uid", "integer#int64", true, "user id"),
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("POST", "/v1/user/following/{uid}", "follow user").
			Tags("Follow").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int64", true, "user id")).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/user/following/{uid}", "unfollow user").
			Tags("Follow").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int64", true, "user id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type FollowController struct {
	config        *config.Config
	jwtService    *service.JwtService
	userService   *service.UserService
	followService *service.FollowService
	common        *CommonController
}

func NewFollowController() *FollowController {
	return &FollowController{
		config:        xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:    xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		userService:   xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		followService: xdi.GetByNameForce(sn.SFollowService).(*service.FollowService),
		common:        xdi.GetByNameForce(sn.SCommonController).(*CommonController),
	}
}

// GET /v1/user/:uid/follower
func (s *FollowController) QueryFollowers(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, s.config)

	users, total, err := s.followService.QueryFollowers(uid, pp)
	if err != nil {
		return result.Error(exception.GetFollowerListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := s.jwtService.GetContextUser(c)
	extras, err := s.common.getUserExtras(c, authUser, users)
	if err != nil {
		return result.Error(exception.GetFollowerListError).SetError(err, c)
	}

	res := dto.BuildUserDtos(users)
	for idx, user := range res {
		user.Extra = extras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// GET /v1/user/:uid/following
func (s *FollowController) QueryFollowings(c *gin.Context) *result.Result {
	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}
	pp := param.BindPageOrder(c, s.config)

	users, total, err := s.followService.QueryFollowings(uid, pp)
	if err != nil {
		return result.Error(exception.GetFollowingListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := s.jwtService.GetContextUser(c)
	extras, err := s.common.getUserExtras(c, authUser, users)
	if err != nil {
		return result.Error(exception.GetFollowingListError).SetError(err, c)
	}

	res := dto.BuildUserDtos(users)
	for idx, user := range res {
		user.Extra = extras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// POST /v1/user/following/:uid
func (s *FollowController) FollowUser(c *gin.Context) *result.Result {
	user := s.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	if user.Uid == uid {
		return result.Error(exception.FollowSelfError)
	}

	status, err := s.followService.FollowUser(user.Uid, uid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbExisted {
		return result.Error(exception.AlreadyFollowingError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.FollowError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/user/following/:uid
func (s *FollowController) UnfollowUser(c *gin.Context) *result.Result {
	user := s.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	status, err := s.followService.UnfollowUser(user.Uid, uid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbTagA {
		return result.Error(exception.NotFollowYetError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UnfollowError).SetError(err, c)
	}

	return result.Ok()
}
