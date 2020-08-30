package profile

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
)

func addPropertyMappers() {
	xproperty.AddMapper(xproperty.NewMapper(&dto.VideoDto{}, &po.Video{}, map[string]*xproperty.PropertyMapperValue{
		"vid":         xproperty.NewValue(false, "vid"),
		"title":       xproperty.NewValue(false, "title"),
		"description": xproperty.NewValue(false, "description"),
		"video_url":   xproperty.NewValue(false, "video_url"),
		"cover_url":   xproperty.NewValue(false, "cover_url"),
		"upload_time": xproperty.NewValue(false, "create_at"),
		"update_time": xproperty.NewValue(false, "update_at"),
		"author_uid":  xproperty.NewValue(false, "author_uid"),
	}))
}
