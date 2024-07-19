package driver

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Println("db:", db)

	if err != nil {
		log.Fatalf("Error when attemping to connect to the database:", err)
		return nil, err
	}

	return db, err
}
