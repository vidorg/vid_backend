package profile

import (
	"github.com/Aoi-hosizora/ahlib/xentity"
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/src/config"
)

func CreateEntityMappers(config *config.ServerConfig) *xentity.EntityMappers {
	mappers := xentity.NewEntityMappers()
	addDtoMappers(config, mappers)
	addParamMappers(config, mappers)
	return mappers
}

func CreatePropertyMappers() *xproperty.PropertyMappers {
	mappers := xproperty.NewPropertyMappers()
	addPropertyMappers(mappers)
	return mappers
}
