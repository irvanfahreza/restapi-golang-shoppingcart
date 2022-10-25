package database

import (
	"log"
	"os"

	// "gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ilmudata/project-golang/models"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance
var Db *gorm.DB

func InitDb() *gorm.DB { // OOP constructor
	Db = ConnectDb()
	return Db
}

func ConnectDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("project-golang.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}

	log.Println("Connected Successfully to Database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	db.AutoMigrate(&models.Product{}, &models.Order{})

	Database = DbInstance{
		Db: db,
	}
	return db
}
