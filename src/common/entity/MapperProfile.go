package entity

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/config"
)

func CreateEntityMappers(config *config.ServerConfig) *xmapper.EntityMappers {
	mappers := xmapper.NewEntityMappers()

	addDtoMappers(config, mappers)
	addParamMappers(config, mappers)

	return mappers
}
