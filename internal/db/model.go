package db

import "gorm.io/gorm"

// type for database connection values
type DBModel struct {
	DB *gorm.DB
}

// Models is the wrapper for ALL modules
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *gorm.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}
