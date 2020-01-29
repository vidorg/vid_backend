package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgorm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
)

const (
	DefaultDeleteAtTimeStamp = "2000-01-01 00:00:00"
)

func SetupDBConn(cfg *config.MySqlConfig) *gorm.DB {
	dbParams := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Name, cfg.Charset,
	)
	db, err := gorm.Open("mysql", dbParams)
	if err != nil {
		log.Fatalln("Failed to connect mysql:", err)
	}
	// Change default deletedAt field behavior
	xgorm.HookDeleteAtField(db, DefaultDeleteAtTimeStamp)

	db.LogMode(cfg.IsLog)
	db.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tbl_" + defaultTableName
	}

	authMigrate := func(value interface{}) {
		rdb := db.AutoMigrate(value)
		if rdb.Error != nil {
			log.Fatalln("Failed to auto migrate model:", rdb.Error)
		}
	}
	authMigrate(&po.User{})
	authMigrate(&po.Account{})
	authMigrate(&po.Video{})

	return db
}
