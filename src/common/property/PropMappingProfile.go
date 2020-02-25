package property

import (
	"github.com/vidorg/vid_backend/src/model/dto"
	"github.com/vidorg/vid_backend/src/model/po"
	"reflect"
)

type PropMappingProfile struct {
	Mappings []*PropMapping
}

func CreatePropMappingProfile() *PropMappingProfile {
	mappings := make([]*PropMapping, 0)

	mappings = append(mappings, NewPropMapping(&dto.UserDto{}, &po.User{}, map[string]*PropMappingValue{
		"uid":           NewPropMappingValue([]string{"uid"}, false),
		"username":      NewPropMappingValue([]string{"username"}, false),
		"sex":           NewPropMappingValue([]string{"sex"}, false),
		"profile":       NewPropMappingValue([]string{"profile"}, false),
		"avatar_url":    NewPropMappingValue([]string{"avatar_url"}, false),
		"birth_time":    NewPropMappingValue([]string{"birth_time"}, false),
		"authority":     NewPropMappingValue([]string{"authority"}, false),
		"phone_number":  NewPropMappingValue([]string{"phone_number"}, false),
		"register_time": NewPropMappingValue([]string{"register_time"}, false),
	}))

	mappings = append(mappings, NewPropMapping(&dto.VideoDto{}, &po.Video{}, map[string]*PropMappingValue{
		"vid":         NewPropMappingValue([]string{"vid"}, false),
		"title":       NewPropMappingValue([]string{"title"}, false),
		"description": NewPropMappingValue([]string{"description"}, false),
		"video_url":   NewPropMappingValue([]string{"video_url"}, false),
		"cover_url":   NewPropMappingValue([]string{"cover_url"}, false),
		"upload_time": NewPropMappingValue([]string{"create_at"}, false),
		"update_time": NewPropMappingValue([]string{"update_at"}, false),
		"author_uid":  NewPropMappingValue([]string{"author_uid"}, false),
	}))

	return &PropMappingProfile{Mappings: mappings}
}

func (p *PropMappingProfile) GetPropertyMapping(dtoModel interface{}, poModel interface{}) *PropMapping {
	for _, m := range p.Mappings {
		if reflect.TypeOf(m.DtoModel) == reflect.TypeOf(dtoModel) && reflect.TypeOf(m.PoModel) == reflect.TypeOf(poModel) {
			return m
		}
	}
	return NewPropMapping(dtoModel, poModel, map[string]*PropMappingValue{})
}
