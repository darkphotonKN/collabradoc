package db

import (
	"log"

	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// DBCon is the global database connection handle
	DBCon *gorm.DB
)

// Init initializes the database connection
func Init(dsn string) {
	var err error

	DBCon, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Perform Migrations
	err = DBCon.AutoMigrate(&model.User{}, &model.Document{}, &model.Comment{}, &model.LiveSession{}, &model.Rating{})

	if err != nil {
		log.Fatal("DB could not be connected to.")
	}

}
