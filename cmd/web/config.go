package main

import (
	"fmt"
	"os"
)

func (app *application) setDSN() {
	loadDBVars := dbVars{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		loadDBVars.DBHost, loadDBVars.DBUser, loadDBVars.DBPassword, loadDBVars.DBName, loadDBVars.DBPort)

	app.config.db.dsn = dsn
}
