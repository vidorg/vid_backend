package controller

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"github.com/vidorg/vid_backend/src/service"
)

var (
	adNeedSubscribeCount = goapidoc.NewQueryParam("need_subscribe_count", "boolean", false, "need subscribe count (user)")
	adNeedIsSubscribe    = goapidoc.NewQueryParam("need_is_subscribe", "boolean", false, "need is subscribe (user)")
	adNeedVideoCount     = goapidoc.NewQueryParam("need_video_count", "boolean", false, "need video count (user)")
	adNeedAuthor         = goapidoc.NewQueryParam("need_author", "boolean", false, "need video author (video)")
)

type CommonController struct {
	subscribeService *service.SubscribeService
	videoService     *service.VideoService
	userService      *service.UserService
}

func NewCommonController() *CommonController {
	return &CommonController{
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
		videoService:     xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
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
			extras[idx].Subscribers = &cnt[0]
			extras[idx].Subscribings = &cnt[1]
		}
	}

	// need_is_subscribe
	if param.BindQueryBool(c, "need_is_subscribe") && authUser != nil {
		arr, err := cmn.subscribeService.CheckSubscribeByUids(authUser.Uid, uids)
		if err != nil {
			return nil, err
		}
		for idx, is := range arr {
			extras[idx].IsSubscriber = &is[0]
			extras[idx].IsSubscribing = &is[1]
		}
	}

	// need_video_count
	if param.BindQueryBool(c, "need_video_count") && authUser != nil {
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
func (cmn *CommonController) getVideosAuthor(c *gin.Context, videos []*po.Video) ([]*po.User, error) {
	authors := make([]*po.User, len(videos))
	uids := make([]uint64, len(videos))
	for idx, video := range videos {
		uids[idx] = video.AuthorUid
	}

	if param.BindQueryBool(c, "need_author") {

		var err error
		authors, err = cmn.userService.QueryByUids(uids)
		if err != nil {
			return nil, err
		}
	}

	return authors, nil
}
