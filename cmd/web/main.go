package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darkphotonKN/collabradoc/internal/db"
	"github.com/darkphotonKN/collabradoc/internal/document"
	"github.com/darkphotonKN/collabradoc/internal/user"
	"github.com/joho/godotenv"
)

type dbVars struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

type dbCfg struct {
	dsn string
}

type config struct {
	port string
	db   dbCfg
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
}

// Set Up Server
func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	fmt.Println("Server running on", app.config.port)
	return srv.ListenAndServe()
}

func main() {
	// Load Environmental Variables
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Setup Config
	config := config{
		port: os.Getenv("APP_PORT"),
	}

	app := &application{
		config: config,
	}

	// set up DSN
	app.setDSN()

	// Connecting to DB
	db.Init(app.config.db.dsn)

	// Perform Migrations
	err = db.DBCon.AutoMigrate(&user.User{}, &document.Document{})

	if err != nil {
		log.Fatal("DB could not be connected to.")
	}

	fmt.Println("DB connected.")

	if err != nil {
		log.Fatalf("Could not initialize DB table products.")
	}

	// Start Server
	app.serve()
}
