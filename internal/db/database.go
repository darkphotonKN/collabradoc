package db

import (
	"log"

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

}
