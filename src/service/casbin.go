package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/casbin/casbin/v2"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type CasbinService struct {
	enforcer *casbin.Enforcer
}

func NewCasbinService() *CasbinService {
	srv := &CasbinService{
		enforcer: xdi.GetByNameForce(sn.SEnforcer).(*casbin.Enforcer),
	}
	return srv
}

func (c *CasbinService) Enforce(sub string, obj string, act string) (bool, error) {
	err := c.enforcer.LoadPolicy()
	if err != nil {
		return false, err
	}

	return c.enforcer.Enforce(sub, obj, act)
}
