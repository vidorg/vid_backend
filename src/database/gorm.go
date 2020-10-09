package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/sirupsen/logrus"
	"github.com/vidorg/vid_backend/lib/xgorm"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func NewMySQLDB() (*gorm.DB, error) {
	cfg := xdi.GetByNameForce(sn.SConfig).(*config.Config)
	log := xdi.GetByNameForce(sn.SLogger).(*logrus.Logger)
	mcfg := cfg.MySQL

	logLevel := logger.Warn
	if cfg.Meta.RunMode == "debug" {
		logLevel = logger.Info
	}

	params := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", mcfg.User, mcfg.Password, mcfg.Host, mcfg.Port, mcfg.Name, mcfg.Charset)
	db, err := gorm.Open(mysql.Open(params), &gorm.Config{
		Logger: xgorm.NewGormLogger(log, logger.Config{LogLevel: logLevel}),
		NamingStrategy: &schema.NamingStrategy{
			TablePrefix:   "tbl_",
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(int(mcfg.MaxIdle))
	sqlDB.SetMaxOpenConns(int(mcfg.MaxActive))
	sqlDB.SetConnMaxLifetime(time.Duration(mcfg.MaxLifetime) * time.Second)

	return db, nil
}

func migrate(db *gorm.DB) error {
	for _, err := range []error{
		db.SetupJoinTable(&po.User{}, "Followings", &po.Follow{}),
		db.SetupJoinTable(&po.User{}, "Followers", &po.Follow{}),
		db.SetupJoinTable(&po.User{}, "Subscribings", &po.Subscribe{}),
		db.SetupJoinTable(&po.Channel{}, "Subscribers", &po.Subscribe{}),
		db.SetupJoinTable(&po.User{}, "Favorites", &po.Favorite{}),
		db.SetupJoinTable(&po.Video{}, "Favoreds", &po.Favorite{}),
	} {
		if err != nil {
			return err
		}
	}

	for _, val := range []interface{}{
		&po.User{}, &po.Account{}, &po.Video{}, &po.Channel{},
	} {
		if err := db.AutoMigrate(val); err != nil {
			return err
		}
	}

	return nil
}
