package database

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

func NewCasbinEnforcer() (*casbin.Enforcer, error) {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	db := xdi.GetByNameForce(sn.SGorm).(*gorm.DB)

	adapter, err := gormadapter.NewAdapterByDBUsePrefix(db, "tbl_")
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(cfg.Casbin.ConfigPath, adapter)
	if err != nil {
		return nil, err
	}

	enforcer.EnableLog(false)
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
