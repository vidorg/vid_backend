package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v2"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type CasbinService struct {
	config     *config.Config
	db         *gorm.DB
	adapter    *gormadapter.Adapter
	jwtService *JwtService
}

func NewCasbinService() *CasbinService {
	srv := &CasbinService{
		config:     xdi.GetByNameForce(sn.SConfig).(*config.Config),
		db:         xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		adapter:    xdi.GetByNameForce(sn.SGormAdapter).(*gormadapter.Adapter),
		jwtService: xdi.GetByNameForce(sn.SJwtService).(*JwtService),
	}
	return srv
}

func (c *CasbinService) GetEnforcer() (*casbin.Enforcer, error) {
	enforcer, err := casbin.NewEnforcer(c.config.Casbin.ConfigPath, c.adapter)
	if err != nil {
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

func (c *CasbinService) Enforce(sub string, obj string, act string) (bool, error) {
	enforcer, err := c.GetEnforcer()
	if err != nil {
		return false, nil
	}

	return enforcer.Enforce(sub, obj, act)
}

func (c *CasbinService) GetRoles() ([]string, bool) {
	enforcer, err := c.GetEnforcer()
	if err != nil {
		return nil, false
	}

	return enforcer.GetAllRoles(), true
}

func (c *CasbinService) GetPolicies(limit int32, page int32) (total int32, policies []*po.Policy) {
	total = 0
	policies = make([]*po.Policy, 0)
	c.db.Table("tbl_casbin_rule").Count(&total)
	c.db.Table("tbl_casbin_rule").Limit(limit).Offset((page - 1) * limit).Find(&policies)

	return total, policies
}

func (c *CasbinService) AddPolicy(sub string, obj string, act string) xstatus.DbStatus {
	enforcer, err := c.GetEnforcer()
	if err != nil {
		return xstatus.DbFailed
	}

	ok, err := enforcer.AddPolicy(sub, obj, act)
	if !ok {
		return xstatus.DbExisted
	} else if err != nil {
		return xstatus.DbFailed
	}
	return xstatus.DbSuccess
}

func (c *CasbinService) DeletePolicy(sub string, obj string, act string) xstatus.DbStatus {
	enforcer, err := c.GetEnforcer()
	if err != nil {
		return xstatus.DbFailed
	}

	ok, err := enforcer.RemovePolicy(sub, obj, act)
	if !ok {
		return xstatus.DbNotFound
	} else if err != nil {
		return xstatus.DbFailed
	}
	return xstatus.DbSuccess
}
