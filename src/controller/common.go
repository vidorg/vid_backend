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
	userService      *service.UserService
	followService    *service.FollowService
	channelService   *service.ChannelService
	subscribeService *service.SubscribeService
	videoService     *service.VideoService
	favoriteService  *service.FavoriteService
}

func NewCommonController() *CommonController {
	return &CommonController{
		userService:      xdi.GetByNameForce(sn.SUserService).(*service.UserService),
		followService:    xdi.GetByNameForce(sn.SFollowService).(*service.FollowService),
		channelService:   xdi.GetByNameForce(sn.SChannelService).(*service.ChannelService),
		subscribeService: xdi.GetByNameForce(sn.SSubscribeService).(*service.SubscribeService),
		videoService:     xdi.GetByNameForce(sn.SVideoService).(*service.VideoService),
		favoriteService:  xdi.GetByNameForce(sn.SFavoriteService).(*service.FavoriteService),
	}
}

// Query []*dto.UserExtraDto from []*po.User.
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
		arr, err := cmn.subscribeService.QuerySubscribingCount(uids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Subscribings = &cnt
		}
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

// Query []*po.User and []*dto.UserExtraDto from []*po.Channel.
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

// Query []*dto.ChannelExtraDto from []*po.Channel.
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
		arr, err := cmn.subscribeService.QuerySubscriberCount(cids)
		if err != nil {
			return nil, err
		}
		for idx, cnt := range arr {
			cnt := cnt
			extras[idx].Subscribers = &cnt
		}
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
			arr, err := cmn.subscribeService.CheckSubscribe(authUser.Uid, cids)
			if err != nil {
				return nil, err
			}
			for idx, is := range arr {
				is := is
				extras[idx].IsSubscribed = &is
			}
		} else {
			for idx := range extras {
				extras[idx].IsSubscribed = xpointer.BoolPtr(false)
			}
		}
	}

	return extras, nil
}

// Query []*po.Channel and []*dto.ChannelExtraDto from []*po.Video.
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

// Query []*dto.VideoExtraDto from []*po.Video.
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

// PreLoad []*dto.UserDto from []*po.User.
func (cmn *CommonController) PreLoadUsers(c *gin.Context, authUser *po.User, users []*po.User, out []*dto.UserDto) error {
	extras, err := cmn.getUserExtras(c, authUser, users)
	if err != nil {
		return err
	}

	for idx, user := range out {
		user.Extra = extras[idx]
	}

	return nil
}

// PreLoad []*dto.ChannelDto from []*po.Channel.
func (cmn *CommonController) PreLoadChannels(c *gin.Context, authUser *po.User, channels []*po.Channel, out []*dto.ChannelDto) error {
	authors, userExtras, err := cmn.getChannelAuthors(c, authUser, channels)
	if err != nil {
		return err
	}
	channelExtras, err := cmn.getChannelExtras(c, authUser, channels)
	if err != nil {
		return err
	}

	for idx, channel := range out {
		channel.Author = dto.BuildUserDto(authors[idx])
		if channel.Author != nil {
			channel.Author.Extra = userExtras[idx]
		}
		channel.Extra = channelExtras[idx]
	}

	return nil
}

// PreLoad []*dto.VideoDto from []*po.Video.
func (cmn *CommonController) PreLoadVideos(c *gin.Context, authUser *po.User, videos []*po.Video, out []*dto.VideoDto) error {
	channels, channelExtras, err := cmn.getVideoChannels(c, authUser, videos)
	if err != nil {
		return err
	}
	authors, userExtras := make([]*po.User, 0), make([]*dto.UserExtraDto, 0)
	if len(channels) != 0 && channels[0] != nil {
		authors, userExtras, err = cmn.getChannelAuthors(c, authUser, channels)
		if err != nil {
			return err
		}
	}
	videoExtras, err := cmn.getVideoExtras(c, authUser, videos)
	if err != nil {
		return err
	}

	for idx, video := range out {
		video.Channel = dto.BuildChannelDto(channels[idx])
		if video.Channel != nil {
			video.Channel.Extra = channelExtras[idx]
			video.Channel.Author = dto.BuildUserDto(authors[idx])
			if video.Channel.Author != nil {
				video.Channel.Author.Extra = userExtras[idx]
			}
		}
		video.Extra = videoExtras[idx]
	}

	return nil
}
