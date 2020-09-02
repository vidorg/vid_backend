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
		goapidoc.NewRoutePath("GET", "/v1/block/user", "query user blockings").
			Tags("Block").
			Securities("Jwt").
			Params(
				param.ADPage, param.ADLimit, param.ADOrder,
				_adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedBlockCount, _adNeedVideoCount,
			).
			Responses(goapidoc.NewResponse(200, "_Result<_Page<UserDto>>")),

		goapidoc.NewRoutePath("POST", "/v1/block/user/{uid}", "block user").
			Tags("Block").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int64", true, "user id")).
			Responses(goapidoc.NewResponse(200, "Result")),

		goapidoc.NewRoutePath("DELETE", "/v1/block/user/{uid}", "unblock user").
			Tags("Block").
			Securities("Jwt").
			Params(goapidoc.NewPathParam("uid", "integer#int64", true, "user id")).
			Responses(goapidoc.NewResponse(200, "Result")),
	)
}

type BlockController struct {
	config       *config.Config
	jwtService   *service.JwtService
	userService  *service.UserService
	blockService *service.BlockService
	common       *CommonController
}

func NewBlockController() *BlockController {
	return &BlockController{
		config:       xdi.GetByNameForce(sn.SConfig).(*config.Config),
		jwtService:   xdi.GetByNameForce(sn.SJwtService).(*service.JwtService),
		userService:  xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		blockService: xdi.GetByNameForce(sn.SBlockService).(*service.BlockService),
		common:       xdi.GetByNameForce(sn.SCommonController).(*CommonController),
	}
}

// GET /v1/block/user
func (b *BlockController) QueryBlockings(c *gin.Context) *result.Result {
	user := b.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}
	pp := param.BindPageOrder(c, b.config)

	users, total, err := b.blockService.QueryBlockings(user.Uid, pp)
	if err != nil {
		return result.Error(exception.GetBlockingListError).SetError(err, c)
	} else if users == nil {
		return result.Error(exception.UserNotFoundError)
	}

	authUser := b.jwtService.GetContextUser(c)
	extras, err := b.common.getUsersExtra(c, authUser, users)
	if err != nil {
		return result.Error(exception.QueryUserError).SetError(err, c)
	}

	res := dto.BuildUserDtos(users)
	for idx, user := range res {
		user.Extra = extras[idx]
	}
	return result.Ok().SetPage(pp.Page, pp.Limit, total, res)
}

// POST /v1/block/user/:uid
func (b *BlockController) BlockUser(c *gin.Context) *result.Result {
	user := b.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	if user.Uid == uid {
		return result.Error(exception.BlockSelfError)
	}

	status, err := b.blockService.InsertBlocking(user.Uid, uid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbExisted {
		return result.Error(exception.AlreadyBlockingError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.BlockError).SetError(err, c)
	}

	return result.Ok()
}

// DELETE /v1/block/user/:uid
func (b *BlockController) UnblockUser(c *gin.Context) *result.Result {
	user := b.jwtService.GetContextUser(c)
	if user == nil {
		return nil
	}

	uid, err := param.BindRouteId(c, "uid")
	if err != nil {
		return result.Error(exception.RequestParamError).SetError(err, c)
	}

	status, err := b.blockService.DeleteBlocking(user.Uid, uid)
	if status == xstatus.DbNotFound {
		return result.Error(exception.UserNotFoundError)
	} else if status == xstatus.DbTagA {
		return result.Error(exception.NotBlockYetError)
	} else if status == xstatus.DbFailed {
		return result.Error(exception.UnblockError).SetError(err, c)
	}

	return result.Ok()
}
