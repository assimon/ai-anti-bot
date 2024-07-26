package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var Gdb *gorm.DB

func init() {
	var err error
	Gdb, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("open database err:", err)
	}
	err = Gdb.AutoMigrate(UserInfo{}, Advertise{})
	if err != nil {
		log.Fatal("database AutoMigrate err:", err)
	}
	database, _ := Gdb.DB()
	database.SetMaxOpenConns(1)
	err = database.Ping()
	if err != nil {
		log.Fatal("database ping err:", err)
	}
}
