package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xpointer"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

var (
	_adNeedSubscribeCount = goapidoc.NewQueryParam("need_subscribe_count", "boolean", false, "need subscribe count (user)")
	_adNeedIsSubscribe    = goapidoc.NewQueryParam("need_is_subscribe", "boolean", false, "need is subscribe (user)")
	_adNeedIsBlock        = goapidoc.NewQueryParam("need_is_block", "boolean", false, "need block count (user)")
	_adNeedVideoCount     = goapidoc.NewQueryParam("need_video_count", "boolean", false, "need video count (user)")
	_adNeedAuthor         = goapidoc.NewQueryParam("need_author", "boolean", false, "need video author (video)")
)

type CommonController struct {
	subscribeService *service.SubscribeService
	videoService     *service.VideoService
	blockService     *service.BlockService
	userService      *service.UserService
}

func NewCommonController() *CommonController {
	return &CommonController{
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
		videoService:     xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
		blockService:     xdi.GetByNameForce(sn.SBlockService).(*service.BlockService),
		userService:      xdi.GetByNameForce(sn.SUserService).(*service.UserService),
	}
}

// Get dto.UserExtraDto for user list.
func (cmn *CommonController) getUsersExtra(c *gin.Context, authUser *po.User, users []*po.User) ([]*dto.UserExtraDto, error) {
	extras := make([]*dto.UserExtraDto, len(users))
	for idx := range extras {
		extras[idx] = &dto.UserExtraDto{}
	}
	uids := make([]uint64, len(users))
	for idx, user := range users {
		uids[idx] = user.Uid
	}

	// need_subscribe_count
	if param.BindQueryBool(c, "need_subscribe_count") {
		arr, err := cmn.subscribeService.QueryCountByUids(uids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			extras[idx].Subscribings = &cnt[0]
			extras[idx].Subscribers = &cnt[1]
		}
	}

	// need_is_subscribe
	if param.BindQueryBool(c, "need_is_subscribe") {
		if authUser != nil {
			arr, err := cmn.subscribeService.CheckSubscribes(authUser.Uid, uids)
			if err != nil {
				return nil, err
			}
			for idx, is := range arr {
				extras[idx].IsSubscribing = &is[0]
				extras[idx].IsSubscribed = &is[1]
			}
		} else {
			for idx := range extras {
				extras[idx].IsSubscribing = xpointer.BoolPtr(false)
				extras[idx].IsSubscribed = xpointer.BoolPtr(false)
			}
		}
	}

	// need_is_block
	if param.BindQueryBool(c, "need_is_block") {
		if authUser != nil {
			arr, err := cmn.blockService.CheckBlockings(authUser.Uid, uids)
			if err != nil {
				return nil, err
			}
			for idx, is := range arr {
				is := is
				extras[idx].IsBlocking = &is
			}
		} else {
			for idx := range extras {
				extras[idx].IsBlocking = xpointer.BoolPtr(false)
			}
		}
	}

	// need_video_count
	if param.BindQueryBool(c, "need_video_count") {
		arr, err := cmn.videoService.QueryCountByUids(uids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Videos = &cnt
		}
	}

	return extras, nil
}

// Get po.Video author for video list.
func (cmn *CommonController) getVideosAuthor(c *gin.Context, authUser *po.User, videos []*po.Video) ([]*po.User, []*dto.UserExtraDto, error) {
	authors := make([]*po.User, len(videos))
	extras := make([]*dto.UserExtraDto, len(videos))
	for idx := range extras {
		extras[idx] = &dto.UserExtraDto{}
	}
	uids := make([]uint64, len(videos))
	for idx, video := range videos {
		uids[idx] = video.AuthorUid
	}

	if param.BindQueryBool(c, "need_author") {
		var err error
		authors, err = cmn.userService.QueryByUids(uids)
		if err != nil {
			return nil, nil, err
		}

		extras, err = cmn.getUsersExtra(c, authUser, authors)
		if err != nil {
			return nil, nil, err
		}
	}

	return authors, extras, nil
}
