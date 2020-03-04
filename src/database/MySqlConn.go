package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgorm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/database/helper"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
)

const (
	DefaultDeleteAtTimeStamp = "2000-01-01 00:00:00"
)

func SetupMySqlConn(cfg *config.MySqlConfig) *helper.GormHelper {
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

	autoMigrateModel(db)
	addFullTextIndex(db, cfg)

	return helper.NewGormHelper(db)
}

func autoMigrateModel(db *gorm.DB) {
	autoMigrate := func(value interface{}) {
		rdb := db.AutoMigrate(value)
		if rdb.Error != nil {
			log.Fatalln(rdb.Error)
		}
	}

	autoMigrate(&po.User{})
	autoMigrate(&po.Account{})
	autoMigrate(&po.Video{})
}

func addFullTextIndex(db *gorm.DB, cfg *config.MySqlConfig) {
	checkExecIndex := func(tblName string, idxName string, param string) {
		cnt := 0
		rdb := db.Table("INFORMATION_SCHEMA.STATISTICS").Where("TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?", cfg.Name, tblName, idxName).Count(&cnt)
		if rdb.Error != nil {
			log.Fatalln(rdb.Error)
		}
		if cnt == 0 {
			sql := fmt.Sprintf("CREATE FULLTEXT INDEX `%s` ON `%s` (%s) WITH PARSER `ngram`", idxName, tblName, param)
			rdb := db.Exec(sql)
			if rdb.Error != nil {
				log.Fatalln(rdb.Error)
			}
		}
	}

	checkExecIndex("tbl_user", "idx_username_profile_fulltext", "`username`(100), `profile`(20)")
	checkExecIndex("tbl_video", "idx_title_description_fulltext", "`title`(100), `description`(40)")
}
