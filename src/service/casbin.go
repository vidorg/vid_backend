package service

import (
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v2"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database"
	"github.com/vidorg/vid_backend/src/model/po"
)

type CasbinService struct {
	Config     *config.Config       `di:"~"`
	Logger     *logrus.Logger       `di:"~"`
	Db         *database.GormHelper `di:"~"`
	JwtService *JwtService          `di:"~"`

	Adapter *gormadapter.Adapter `di:"-"`
}

func NewCasbinService(dic *xdi.DiContainer) *CasbinService {
	srv := &CasbinService{}
	dic.MustInject(srv)

	adapter, err := gormadapter.NewAdapterByDBUsePrefix(srv.Db.DB, "tbl_")
	if err != nil {
		panic(err)
	}
	srv.Adapter = adapter
	return srv
}

func (c *CasbinService) GetEnforcer() (*casbin.Enforcer, error) {
	enforcer, err := casbin.NewEnforcer(c.Config.Casbin.ConfigPath, c.Adapter)
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

func (c *CasbinService) GetPolicies(limit int32, page int32) (int32, []*po.Policy) {
	total := 0
	policies := make([]*po.Policy, 0)
	c.Db.Table("tbl_casbin_rule").Count(&total)
	c.Db.Table("tbl_casbin_rule").Limit(limit).Offset((page - 1) * limit).Find(&policies)
	return int32(total), policies
}

func (c *CasbinService) AddPolicy(sub string, obj string, act string) database.DbStatus {
	enforcer, err := c.GetEnforcer()
	if err != nil {
		return database.DbFailed
	}
	ok, err := enforcer.AddPolicy(sub, obj, act)
	if !ok {
		return database.DbExisted
	} else if err != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}

func (c *CasbinService) DeletePolicy(sub string, obj string, act string) database.DbStatus {
	enforcer, err := c.GetEnforcer()
	if err != nil {
		return database.DbFailed
	}
	ok, err := enforcer.RemovePolicy(sub, obj, act)
	if !ok {
		return database.DbNotFound
	} else if err != nil {
		return database.DbFailed
	}
	return database.DbSuccess
}
