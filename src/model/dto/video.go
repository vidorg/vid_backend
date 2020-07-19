package dto

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("VideoDto", "视频信息").
			WithProperties(
				goapidoc.NewProperty("vid", "integer#int32", true, "视频id"),
				goapidoc.NewProperty("title", "string", true, "标题"),
				goapidoc.NewProperty("description", "string", true, "简介").WithAllowEmptyValue(true),
				goapidoc.NewProperty("video_url", "string", true, "资源"),
				goapidoc.NewProperty("cover_url", "string", true, "封面"),
				goapidoc.NewProperty("upload_time", "string#date-time", true, "上传时间"),
				goapidoc.NewProperty("update_time", "string#date-time", true, "修改时间"),
				goapidoc.NewProperty("author", "UserDto", true, "视频作者, null 表示用户不存在"),
			),
	)
}

type VideoDto struct {
	Vid         int32    `json:"vid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	VideoUrl    string   `json:"video_url"` // TODO url
	CoverUrl    string   `json:"cover_url"` // TODO url
	UploadTime  string   `json:"upload_time"`
	UpdateTime  string   `json:"update_time"`
	Author      *UserDto `json:"author"`
}

func BuildVideoDto(video *po.Video) *VideoDto {
	return xentity.MustMap(video, &VideoDto{}).(*VideoDto)
}

func BuildVideoDtos(videos []*po.Video) []*VideoDto {
	return xentity.MustMapSlice(videos, &VideoDto{}).([]*VideoDto)
}
