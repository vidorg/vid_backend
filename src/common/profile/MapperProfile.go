package profile

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/config"
)

func CreateEntityMappers(config *config.ServerConfig) *xmapper.EntityMappers {
	mappers := xmapper.NewEntityMappers()

	mappers = loadDtoProfile(config, mappers)
	mappers = loadParamProfile(config, mappers)

	return mappers
}
