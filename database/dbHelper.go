package database

import (
	"fmt"
	"log"
	"vid/config"
	"vid/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func SetupDBConn(cfg config.Config) {
	dbParams := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbName,
	)
	var err error
	DB, err = gorm.Open("mysql", dbParams)
	if err != nil {
		log.Fatal(2, err)
	}

	DB.LogMode(true)

	DB.AutoMigrate(&models.User{})
}
