package dto

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("VideoDto", "视频信息").
			Properties(
				goapidoc.NewProperty("vid", "integer#int64", true, "video id"),
				goapidoc.NewProperty("title", "string", true, "video title"),
				goapidoc.NewProperty("description", "string", true, "video description"),
				goapidoc.NewProperty("video_url", "string", true, "video source url"),
				goapidoc.NewProperty("cover_url", "string", true, "video cover url"),
				goapidoc.NewProperty("upload_time", "string#date-time", true, "video upload time"),
				goapidoc.NewProperty("author", "UserDto", true, "video author"),
			),
	)
}

type VideoDto struct {
	Vid         uint64   `json:"vid"`         // video id
	Title       string   `json:"title"`       // video title
	Description string   `json:"description"` // video description
	VideoUrl    string   `json:"video_url"`   // video source url (oss)
	CoverUrl    string   `json:"cover_url"`   // video cover url (oss)
	UploadTime  string   `json:"upload_time"` // video upload time
	Author      *UserDto `json:"author"`      // video author
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
		UploadTime:  xtime.NewJsonDateTime(video.CreatedAt).String(),
		Author:      BuildUserDto(video.Author),
	}
}

func BuildVideoDtos(videos []*po.Video) []*VideoDto {
	out := make([]*VideoDto, len(videos))
	for idx, video := range videos {
		out[idx] = BuildVideoDto(video)
	}
	return out
}

func BuildVideoPropertyMapper() xproperty.PropertyDict {
	return xproperty.PropertyDict{
		"vid":         xproperty.NewValue(false, "vid"),
		"title":       xproperty.NewValue(false, "title"),
		"description": xproperty.NewValue(false, "description"),
		"video_url":   xproperty.NewValue(false, "video_url"),
		"cover_url":   xproperty.NewValue(false, "cover_url"),
		"upload_time": xproperty.NewValue(false, "create_at"),
		"author_uid":  xproperty.NewValue(false, "author_uid"),
	}
}
