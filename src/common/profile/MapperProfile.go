package profile

import (
	"github.com/Aoi-hosizora/ahlib/xmapper"
	"github.com/vidorg/vid_backend/src/config"
)

func CreateEntityMapper(config *config.ServerConfig) *xmapper.EntityMapper {
	mapper := xmapper.NewEntityMapper()

	mapper = loadDtoProfile(config, mapper)
	mapper = loadParamProfile(config, mapper)

	return mapper
}
