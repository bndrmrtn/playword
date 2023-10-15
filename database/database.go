package database

import (
	"fmt"
	"github.com/bndrmrtn/playword/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	dsn := "root:@tcp(127.0.0.1:3306)/go_playword?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v\n", err.Error()))
	}

	log.Println("Connected to the database successfully!")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations")
	err2 := db.AutoMigrate(&models.Word{})
	if err2 != nil {
		panic("Failed to migrate: " + err2.Error())
	}

	Database = DbInstance{Db: db}
	fillDb()
}

func fillDb() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	mySqlString, err2 := os.ReadFile(currentDir + "/database/db.sql")

	if err2 != nil {
		panic(err2.Error())
	}

	var count int64
	Database.Db.Model(models.Word{}).Count(&count)

	if count == 0 {
		Database.Db.Exec(string(mySqlString))
	}
}
