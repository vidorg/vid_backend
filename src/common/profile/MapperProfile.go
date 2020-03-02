package profile

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/src/config"
)

func CreateEntityMappers(config *config.ServerConfig) *xmapper.EntityMappers {
	mappers := xmapper.NewEntityMappers()

	addDtoMappers(config, mappers)
	addParamMappers(config, mappers)

	return mappers
}

func CreatePropertyMappers() *xproperty.PropertyMappers {
	mappers := xproperty.NewPropertyMappers()

	addPropMappers(mappers)

	return mappers
}
