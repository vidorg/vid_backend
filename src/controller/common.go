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
	// user
	_adNeedFollowCount      = goapidoc.NewQueryParam("need_follow_count", "boolean", false, "need follow count (user)")
	_adNeedChannelCount     = goapidoc.NewQueryParam("need_channel_count", "boolean", false, "need channel count (user)")
	_adNeedSubscribingCount = goapidoc.NewQueryParam("need_subscribing_count", "boolean", false, "need subscribing count (user)")
	_adNeedFavoriteCount    = goapidoc.NewQueryParam("need_favorite_count", "boolean", false, "need favorite count (user)")
	_adNeedIsFollow         = goapidoc.NewQueryParam("need_is_follow", "boolean", false, "need is follow (user)")

	// channel
	_adNeedChannelAuthor   = goapidoc.NewQueryParam("need_channel_author", "boolean", false, "need channel author (channel)")
	_adNeedSubscriberCount = goapidoc.NewQueryParam("need_subscriber_count", "boolean", false, "need subscriber count (channel)")
	_adNeedVideoCount      = goapidoc.NewQueryParam("need_video_count", "boolean", false, "need video count (channel)")
	_adNeedIsSubscribed    = goapidoc.NewQueryParam("need_is_subscribed", "boolean", false, "need is subscribed (channel)")

	// video
	_adNeedVideoChannel = goapidoc.NewQueryParam("need_video_channel", "boolean", false, "need video channel (video)")
	_adNeedFavoredCount = goapidoc.NewQueryParam("need_favored_count", "boolean", false, "need favored user count (video)")
	_adNeedIsFavorite   = goapidoc.NewQueryParam("need_is_favorite", "boolean", false, "need is favorite (video)")
)

// noinspection GoUnusedGlobalVariable
var (
	_adUserQueries    = []*goapidoc.Param{_adNeedFollowCount, _adNeedChannelCount, _adNeedSubscribingCount, _adNeedFavoriteCount, _adNeedIsFollow}
	_adChannelQueries = []*goapidoc.Param{_adNeedChannelAuthor, _adNeedSubscriberCount, _adNeedVideoCount, _adNeedIsSubscribed}
	_adVideoQueries   = []*goapidoc.Param{_adNeedVideoChannel, _adNeedFavoredCount, _adNeedIsFavorite}
)

type CommonController struct {
	userService     *service.UserService
	followService   *service.FollowService
	channelService  *service.ChannelService
	videoService    *service.VideoService
	favoriteService *service.FavoriteService
}

func NewCommonController() *CommonController {
	return &CommonController{
		userService:     xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		followService:   xdi.GetByNameForce(sn.SFollowService).(*service.FollowService),
		channelService:  xdi.GetByNameForce(sn.SChannelService).(*service.ChannelService),
		videoService:    xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
		favoriteService: xdi.GetByNameForce(sn.SFavoriteService).(*service.FavoriteService),
	}
}

func (cmn *CommonController) getUserExtras(c *gin.Context, authUser *po.User, users []*po.User) ([]*dto.UserExtraDto, error) {
	extras := make([]*dto.UserExtraDto, len(users))
	for idx := range extras {
		extras[idx] = &dto.UserExtraDto{}
	}
	uids := make([]uint64, len(users))
	for idx, user := range users {
		uids[idx] = user.Uid
	}

	// need_follow_count
	if param.BindQueryBool(c, "need_follow_count") {
		arr, err := cmn.followService.QueryFollowCount(uids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			extras[idx].Followings = &cnt[0]
			extras[idx].Followers = &cnt[1]
		}
	}

	// need_channel_count
	if param.BindQueryBool(c, "need_channel_count") {
		arr, err := cmn.channelService.QueryCountByUids(uids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Channels = &cnt
		}
	}

	// need_subscribing_count
	if param.BindQueryBool(c, "need_subscribing_count") {
		// TODO
	}

	// need_favorite_count
	if param.BindQueryBool(c, "need_favorite_count") {
		arr, err := cmn.favoriteService.QueryFavoriteCount(uids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Favorites = &cnt
		}
	}

	// need_is_follow
	if param.BindQueryBool(c, "need_is_follow") {
		if authUser != nil {
			arr, err := cmn.followService.CheckFollow(authUser.Uid, uids)
			if err != nil {
				return nil, err
			}
			for idx, is := range arr {
				extras[idx].IsFollowing = &is[0]
				extras[idx].IsFollowed = &is[1]
			}
		} else {
			for idx := range extras {
				extras[idx].IsFollowing = xpointer.BoolPtr(false)
				extras[idx].IsFollowed = xpointer.BoolPtr(false)
			}
		}
	}

	return extras, nil
}

func (cmn *CommonController) getChannelAuthors(c *gin.Context, authUser *po.User, channels []*po.Channel) ([]*po.User, []*dto.UserExtraDto, error) {
	authors := make([]*po.User, len(channels))
	extras := make([]*dto.UserExtraDto, len(channels))
	for idx := range extras {
		extras[idx] = &dto.UserExtraDto{}
	}
	uids := make([]uint64, len(channels))
	for idx, video := range channels {
		uids[idx] = video.AuthorUid
	}

	// need_channel_author
	if param.BindQueryBool(c, "need_channel_author") {
		var err error
		authors, err = cmn.userService.QueryByUids(uids)
		if err != nil {
			return nil, nil, err
		}
		extras, err = cmn.getUserExtras(c, authUser, authors)
		if err != nil {
			return nil, nil, err
		}
	}

	return authors, extras, nil
}

func (cmn *CommonController) getChannelExtras(c *gin.Context, authUser *po.User, channels []*po.Channel) ([]*dto.ChannelExtraDto, error) {
	extras := make([]*dto.ChannelExtraDto, len(channels))
	for idx := range extras {
		extras[idx] = &dto.ChannelExtraDto{}
	}
	cids := make([]uint64, len(channels))
	for idx, channel := range channels {
		cids[idx] = channel.Cid
	}

	// need_subscriber_count
	if param.BindQueryBool(c, "need_subscriber_count") {
		// TODO
	}

	// need_video_count
	if param.BindQueryBool(c, "need_video_count") {
		arr, err := cmn.videoService.QueryCountByCids(cids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Videos = &cnt
		}
	}

	// need_is_subscribed
	if param.BindQueryBool(c, "need_is_subscribed") {
		if authUser != nil {
			// TODO
		} else {
			for idx := range extras {
				extras[idx].IsSubscribed = xpointer.BoolPtr(false)
			}
		}
	}

	return extras, nil
}

func (cmn *CommonController) getVideoChannels(c *gin.Context, authUser *po.User, videos []*po.Video) ([]*po.Channel, []*dto.ChannelExtraDto, error) {
	channels := make([]*po.Channel, len(videos))
	extras := make([]*dto.ChannelExtraDto, len(videos))
	for idx := range extras {
		extras[idx] = &dto.ChannelExtraDto{}
	}
	cids := make([]uint64, len(videos))
	for idx, video := range videos {
		cids[idx] = video.ChannelCid
	}

	// need_video_channel
	if param.BindQueryBool(c, "need_video_channel") {
		var err error
		channels, err = cmn.channelService.QueryByCids(cids)
		if err != nil {
			return nil, nil, err
		}
		extras, err = cmn.getChannelExtras(c, authUser, channels)
		if err != nil {
			return nil, nil, err
		}
	}

	return channels, extras, nil
}

func (cmn *CommonController) getVideoExtras(c *gin.Context, authUser *po.User, videos []*po.Video) ([]*dto.VideoExtraDto, error) {
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
		arr, err := cmn.favoriteService.QueryFavoredCount(vids)
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
