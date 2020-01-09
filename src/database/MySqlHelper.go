package database

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"log"
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
	authMigrate(&po.PassRecord{})
	authMigrate(&po.Video{})

	// Change default deletedAt field behavior
	db.Callback().Query().Before("gorm:query").Register("new_deleted_at_before_query_callback", newBeforeQueryUpdateCallback)
	db.Callback().RowQuery().Before("gorm:row_query").Register("new_deleted_at_before_row_query_callback", newBeforeQueryUpdateCallback)
	db.Callback().Update().Before("gorm:update").Register("new_deleted_at_before_update_callback", newBeforeQueryUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", newDeleteCallback)

	return db
}

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == 1062
}
