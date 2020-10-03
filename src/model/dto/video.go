package dto

import (
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("VideoDto", "video response").
			Properties(
				goapidoc.NewProperty("vid", "integer#int64", true, "video id"),
				goapidoc.NewProperty("title", "string", true, "video title"),
				goapidoc.NewProperty("description", "string", true, "video description"),
				goapidoc.NewProperty("video_url", "string", true, "video source url"),
				goapidoc.NewProperty("cover_url", "string", true, "video cover url"),
				goapidoc.NewProperty("channel_cid", "integer#int64", true, "video channel id"),
				goapidoc.NewProperty("channel", "UserDto", true, "video channel"),
				goapidoc.NewProperty("upload_time", "string#date-time", true, "video upload time"),
				goapidoc.NewProperty("extra", "VideoExtraDto", true, "video extra information"),
			),

		goapidoc.NewDefinition("VideoExtraDto", "video extra response").
			Properties(
				goapidoc.NewProperty("favoreds", "integer#int32", true, "video favored count"),
				goapidoc.NewProperty("is_favorite", "boolean", true, "is favoring this video"),
			),
	)
}

type VideoDto struct {
	Vid         uint64         `json:"vid"`         // video id
	Title       string         `json:"title"`       // video title
	Description string         `json:"description"` // video description
	VideoUrl    string         `json:"video_url"`   // video source url (oss)
	CoverUrl    string         `json:"cover_url"`   // video cover url (oss)
	ChannelCid  uint64         `json:"channel_cid"` // video channel cid
	Channel     *ChannelDto    `json:"channel"`     // video channel
	UploadTime  string         `json:"upload_time"` // video upload time
	Extra       *VideoExtraDto `json:"extra"`       // video extra information
}

func BuildVideoDto(video *po.Video) *VideoDto {
	if video == nil {
		return nil
	}
	return &VideoDto{
		Vid:         video.Vid,
		Title:       video.Title,
		Description: video.Description,
		VideoUrl:    video.VideoUrl,
		CoverUrl:    video.CoverUrl,
		ChannelCid:  video.ChannelCid,
		Channel:     BuildChannelDto(video.Channel),
		UploadTime:  xtime.NewJsonDateTime(video.CreatedAt).String(),
		Extra:       &VideoExtraDto{},
	}
}

func BuildVideoDtos(videos []*po.Video) []*VideoDto {
	out := make([]*VideoDto, len(videos))
	for idx, video := range videos {
		out[idx] = BuildVideoDto(video)
	}
	return out
}

type VideoExtraDto struct {
	Favoreds   *int32 `json:"favoreds"`    // video favored count
	IsFavorite *bool  `json:"is_favorite"` // is favoring this video
}
