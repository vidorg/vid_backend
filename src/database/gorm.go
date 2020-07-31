package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/ahlib/xdi"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
)

func NewMySQLConn() (*gorm.DB, error) {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	mcfg := cfg.MySQL
	logger := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)

	dbParams := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", mcfg.User, mcfg.Password, mcfg.Host, mcfg.Port, mcfg.Name, mcfg.Charset)
	db, err := gorm.Open("mysql", dbParams)
	if err != nil {
		return nil, err
	}

	db.LogMode(cfg.Meta.RunMode == "debug")
	db.SetLogger(xgorm.NewGormLogrus(logger))
	db.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tbl_" + defaultTableName
	}
	xgorm.HookDeleteAtField(db, xgorm.DefaultDeleteAtTimeStamp)

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	for _, val := range []interface{}{
		&po.User{}, &po.Account{}, &po.Video{},
	} {
		rdb := db.AutoMigrate(val)
		if rdb.Error != nil {
			return rdb.Error
		}
	}
	return nil
}

func NewGormAdapter() (*gormadapter.Adapter, error) {
	db := xdi.GetByNameForce(sn.SGorm).(*gorm.DB)

	adapter, err := gormadapter.NewAdapterByDBUsePrefix(db, "tbl_")
	if err != nil {
		return nil, err
	}
	return adapter, nil
}
