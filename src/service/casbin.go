package service

import (
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/param"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

type CasbinService struct {
	config     *config.Config
	db         *gorm.DB
	enforcer   *casbin.Enforcer
	jwtService *JwtService
	tblName    string
}

func NewCasbinService() *CasbinService {
	return &CasbinService{
		config:     xdi.GetByNameForce(sn.SConfig).(*config.Config),
		db:         xdi.GetByNameForce(sn.SGorm).(*gorm.DB),
		enforcer:   xdi.GetByNameForce(sn.SEnforcer).(*casbin.Enforcer),
		jwtService: xdi.GetByNameForce(sn.SJwtService).(*JwtService),
		tblName:    "tbl_casbin_rule",
	}
}

func (c *CasbinService) Enforce(sub string, obj string, act string) (bool, error) {
	err := c.enforcer.LoadPolicy()
	if err != nil {
		return false, err
	}

	return c.enforcer.Enforce(sub, obj, act)
}

func (c *CasbinService) GetAllRules(pp *param.PageParam) ([]*po.RbacRule, int32, error) {
	total := int32(0)
	rules := make([]*po.RbacRule, 0)
	c.db.Table(c.tblName).Count(&total)
	rdb := xgorm.WithDB(c.db).Pagination(pp.Limit, pp.Page).Table(c.tblName).Order("p_type").Order("v0").Find(&rules)

	if rdb.Error != nil {
		return nil, 0, rdb.Error
	}
	return rules, total, nil
}

func (c *CasbinService) addRule(rule *po.RbacRule) (xstatus.DbStatus, error) {
	rdb := c.db.Table(c.tblName).Where(rule.ToMap()).First(&po.RbacRule{})
	if !rdb.RecordNotFound() {
		return xstatus.DbExisted, nil
	} else if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	}

	rdb = c.db.Table(c.tblName).Create(rule)
	if rdb.RowsAffected == 0 || rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	}
	return xstatus.DbSuccess, nil
}

func (c *CasbinService) removeRule(rule *po.RbacRule) (xstatus.DbStatus, error) {
	ruleMap := rule.ToMap()
	rdb := c.db.Table(c.tblName).Where(ruleMap).First(&po.RbacRule{})
	if rdb.RecordNotFound() {
		return xstatus.DbNotFound, nil
	} else if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	}

	rdb = c.db.Table(c.tblName).Where(ruleMap).Delete(rule)
	if rdb.RowsAffected == 0 || rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	}
	return xstatus.DbSuccess, nil
}

func (c *CasbinService) AddPolicy(sub, obj, act string) (xstatus.DbStatus, error) {
	rule := &po.RbacRule{PType: "p", V0: sub, V1: obj, V2: act}
	return c.addRule(rule)
}

func (c *CasbinService) AddSubject(sub, sub2 string) (xstatus.DbStatus, error) {
	rule := &po.RbacRule{PType: "g", V0: sub, V1: sub2}
	return c.addRule(rule)
}

func (c *CasbinService) RemovePolicy(sub, obj, act string) (xstatus.DbStatus, error) {
	rule := &po.RbacRule{PType: "p", V0: sub, V1: obj, V2: act}
	return c.removeRule(rule)
}

func (c *CasbinService) RemoveSubject(sub, sub2 string) (xstatus.DbStatus, error) {
	rule := &po.RbacRule{PType: "g", V0: sub, V1: sub2}
	return c.removeRule(rule)
}
