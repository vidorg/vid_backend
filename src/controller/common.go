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
	_adNeedVideoCount     = goapidoc.NewQueryParam("need_video_count", "boolean", false, "need video count (user)")
	_adNeedFavoriteCount  = goapidoc.NewQueryParam("need_favorite_count", "boolean", false, "need favorite count (user)")
	_adNeedAuthor         = goapidoc.NewQueryParam("need_author", "boolean", false, "need video author (video)")
	_adNeedFavoredCount   = goapidoc.NewQueryParam("need_favored_count", "boolean", false, "need favored user count (video)")
	_adNeedIsFavorite     = goapidoc.NewQueryParam("need_is_favorite", "boolean", false, "need is favorite (video)")
)

// noinspection GoUnusedGlobalVariable
var (
	_adUserQueries  = []*goapidoc.Param{_adNeedSubscribeCount, _adNeedIsSubscribe, _adNeedVideoCount, _adNeedFavoriteCount}
	_adVideoQueries = []*goapidoc.Param{_adNeedAuthor, _adNeedFavoredCount, _adNeedIsFavorite}
)

type CommonController struct {
	subscribeService *service.SubscribeService
	videoService     *service.VideoService
	userService      *service.UserService
	favoriteService  *service.FavoriteService
}

func NewCommonController() *CommonController {
	return &CommonController{
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
		videoService:     xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
		userService:      xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		favoriteService:  xdi.GetByNameForce(sn.SFavoriteService).(*service.FavoriteService),
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
		arr, err := cmn.subscribeService.QuerySubscribeCount(uids)
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
			arr, err := cmn.subscribeService.CheckSubscribe(authUser.Uid, uids)
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

	// need_favorite_count
	if param.BindQueryBool(c, "need_favorite_count") {
		arr, err := cmn.favoriteService.QueryFavoredCount(uids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Favorites = &cnt
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

// Get dto.UserExtraDto for user list.
func (cmn *CommonController) getVideosExtra(c *gin.Context, authUser *po.User, videos []*po.Video) ([]*dto.VideoExtraDto, error) {
	extras := make([]*dto.VideoExtraDto, len(videos))
	for idx := range extras {
		extras[idx] = &dto.VideoExtraDto{}
	}
	vids := make([]uint64, len(videos))
	for idx, video := range videos {
		vids[idx] = video.Vid
	}

	// need_favored_count
	if param.BindQueryBool(c, "need_favored_count") {
		arr, err := cmn.favoriteService.QueryFavoriteCount(vids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Favoreds = &cnt
		}
	}

	// need_is_favorite
	if param.BindQueryBool(c, "need_is_favorite") {
		if authUser != nil {
			arr, err := cmn.favoriteService.CheckFavorite(authUser.Uid, vids)
			if err != nil {
				return nil, err
			}
			for idx, is := range arr {
				is := is
				extras[idx].IsFavorite = &is
			}
		} else {
			for idx := range videos {
				extras[idx].IsFavorite = xpointer.BoolPtr(false)
			}
		}
	}

	return extras, nil
}
