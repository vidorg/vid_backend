package param

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/model/po"
)

func init() {
	goapidoc.AddDefinitions(
		goapidoc.NewDefinition("VideoParam", "视频请求参数").
			WithProperties(
				goapidoc.NewProperty("title", "string", true, "标题，长度在 [1, 100] 之间"),
				goapidoc.NewProperty("description", "string", true, "简介，长度在 [0, 255] 之间").WithAllowEmptyValue(true),
				goapidoc.NewProperty("cover_url", "string", true, "封面"),
				goapidoc.NewProperty("video_url", "string", true, "资源"),
			),
	)
}

type VideoParam struct {
	Title       string  `form:"title"       json:"title"       binding:"required,min=1,max=100"`
	Description *string `form:"description" json:"description" binding:"required,min=0,max=255"`
	CoverUrl    string  `form:"cover_url"   json:"cover_url"   binding:"required,url"` // TODO url
	VideoUrl    string  `form:"video_url"   json:"video_url"   binding:"required"`     // TODO url
}

func MapVideoParam(param *VideoParam, video *po.Video) {
	xentity.MustMapProp(param, video)
}
