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

var DB *gorm.DB

var PageSize int

func SetupDBConn(cfg *config.DatabaseConfig) {
	dbParams := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Name, cfg.Charset,
	)
	PageSize = cfg.PageSize
	var err error
	DB, err = gorm.Open("mysql", dbParams)
	if err != nil {
		log.Fatal(err)
	}
	DB.LogMode(cfg.IsLog)
	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tbl_" + defaultTableName
	}

	DB.AutoMigrate(&po.User{})
	DB.AutoMigrate(&po.Password{})
	DB.AutoMigrate(&po.Video{})
}

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == 1062
}

func IsNotFoundError(err error) bool {
	return err == nil && gorm.IsRecordNotFoundError(err)
}
