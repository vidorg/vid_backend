package profile

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
)

func addPropertyMappers(mappers *xproperty.PropertyMappers) {
	mappers.AddMapper(xproperty.NewPropertyMapper(&dto.UserDto{}, &po.User{}, map[string]*xproperty.PropertyMapperValue{
		"uid":           xproperty.NewPropertyMapperValue([]string{"uid"}, false),
		"username":      xproperty.NewPropertyMapperValue([]string{"username"}, false),
		"sex":           xproperty.NewPropertyMapperValue([]string{"sex"}, false),
		"profile":       xproperty.NewPropertyMapperValue([]string{"profile"}, false),
		"avatar_url":    xproperty.NewPropertyMapperValue([]string{"avatar_url"}, false),
		"birthday":      xproperty.NewPropertyMapperValue([]string{"birthday"}, false),
		"authority":     xproperty.NewPropertyMapperValue([]string{"authority"}, false),
		"phone_number":  xproperty.NewPropertyMapperValue([]string{"phone_number"}, false),
		"register_time": xproperty.NewPropertyMapperValue([]string{"register_time"}, false),
	}))

	mappers.AddMapper(xproperty.NewPropertyMapper(&dto.VideoDto{}, &po.Video{}, map[string]*xproperty.PropertyMapperValue{
		"vid":         xproperty.NewPropertyMapperValue([]string{"vid"}, false),
		"title":       xproperty.NewPropertyMapperValue([]string{"title"}, false),
		"description": xproperty.NewPropertyMapperValue([]string{"description"}, false),
		"video_url":   xproperty.NewPropertyMapperValue([]string{"video_url"}, false),
		"cover_url":   xproperty.NewPropertyMapperValue([]string{"cover_url"}, false),
		"upload_time": xproperty.NewPropertyMapperValue([]string{"create_at"}, false),
		"update_time": xproperty.NewPropertyMapperValue([]string{"update_at"}, false),
		"author_uid":  xproperty.NewPropertyMapperValue([]string{"author_uid"}, false),
	}))
}
