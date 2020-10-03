package dto

import (
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("ChannelDto", "video response").
			Properties(
				goapidoc.NewProperty("cid", "integer#int64", true, "channel id"),
				goapidoc.NewProperty("name", "string", true, "channel name"),
				goapidoc.NewProperty("description", "string", true, "channel description"),
				goapidoc.NewProperty("cover_url", "string", true, "channel cover url"),
				goapidoc.NewProperty("author_uid", "integer#int64", true, "channel author id"),
				goapidoc.NewProperty("author", "UserDto", true, "channel author"),
				goapidoc.NewProperty("create_time", "string#date-time", true, "channel upload time"),
				goapidoc.NewProperty("extra", "ChannelExtraDto", true, "channel extra information"),
			),

		goapidoc.NewDefinition("ChannelExtraDto", "channel extra response").
			Properties(
				goapidoc.NewProperty("is_subscribed", "integer#int32", true, "channel subscriber users count"),
				goapidoc.NewProperty("is_subscribed", "boolean", true, "this channel is subscribed by the authorized user"),
			),
	)
}

type ChannelDto struct {
	Cid         uint64           `json:"cid"`         // channel id
	Name        string           `json:"name"`        // channel name
	Description string           `json:"description"` // channel description
	CoverUrl    string           `json:"cover_url"`   // channel cover url
	AuthorUid   uint64           `json:"author_uid"`  // channel author id
	Author      *UserDto         `json:"author"`      // channel author
	CreateTime  string           `json:"create_time"` // channel create time
	Extra       *ChannelExtraDto `json:"extra"`       // channel extra information
}

func BuildChannelDto(channel *po.Channel) *ChannelDto {
	if channel == nil {
		return nil
	}
	return &ChannelDto{
		Cid:         channel.Cid,
		Name:        channel.Name,
		Description: channel.Description,
		CoverUrl:    channel.CoverUrl,
		AuthorUid:   channel.AuthorUid,
		Author:      BuildUserDto(channel.Author),
		CreateTime:  xtime.NewJsonDateTime(channel.CreatedAt).String(),
		Extra:       &ChannelExtraDto{},
	}
}

func BuildChannelDtos(channels []*po.Channel) []*ChannelDto {
	out := make([]*ChannelDto, len(channels))
	for idx, channel := range channels {
		out[idx] = BuildChannelDto(channel)
	}
	return out
}

type ChannelExtraDto struct {
	Subscribers  *int32 `json:"subscribers"`   // channel subscriber users count
	IsSubscribed *bool  `json:"is_subscribed"` // this channel is subscribed by the authorized user
}
