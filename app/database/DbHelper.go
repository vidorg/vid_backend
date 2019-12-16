package database

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"vid/app/config"
	"vid/app/model/po"
)

var DB *gorm.DB

func SetupDBConn(cfg config.DatabaseConfig) {
	dbParams := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=%v&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Name,
		cfg.Charset,
	)
	var err error
	DB, err = gorm.Open("mysql", dbParams)
	if err != nil {
		log.Fatal(err)
	}

	DB.LogMode(true)
	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tbl_" + defaultTableName
	}

	DB.AutoMigrate(&po.User{})
	DB.AutoMigrate(&po.PassRecord{})

	// DB.AutoMigrate(&po.Video{})
	// DB.AutoMigrate(&po.Playlist{})
	// DB.AutoMigrate(&po.VideoList{})
}

func IsDuplicateError(err error) bool {
	// https://github.com/jinzhu/gorm/issues/1718
	if err == nil {
		return false
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == 1062
}

func IsNotFoundError(err error) bool {
	return err == nil && gorm.IsRecordNotFoundError(err)
}
