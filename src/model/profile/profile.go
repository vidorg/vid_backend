package profile

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

func BuildEntityMappers() {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	addDtoMappers(cfg)
	addParamMappers(cfg)
}

func BuildPropertyMappers() {
	addPropertyMappers()
}
