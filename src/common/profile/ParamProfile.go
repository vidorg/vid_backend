package profile

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/common/enum"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
)

func loadParamProfile(config *config.ServerConfig, mapper *xmapper.EntityMapper) *xmapper.EntityMapper {
	mapper = mapper.
		CreateMapper(&param.UserParam{}, &po.User{}).
		ForMember("Sex", func(i interface{}) interface{} {
			return enum.ParseSexType(i.(param.UserParam).Sex)
		}).
		ForMember("Profile", func(i interface{}) interface{} {
			return *i.(param.UserParam).Profile
		}).
		ForMember("BirthTime", func(i interface{}) interface{} {
			t, _ := xdatetime.JsonDate{}.Parse(i.(param.UserParam).BirthTime, config.MetaConfig.CurrentLoc)
			return t
		}).
		Build()

	mapper = mapper.
		CreateMapper(&param.VideoParam{}, &po.Video{}).
		ForMember("Description", func(i interface{}) interface{} {
			return *i.(param.VideoParam).Description
		}).Build()

	return mapper
}
