package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darkphotonKN/collabradoc/internal/db"
	"github.com/darkphotonKN/collabradoc/internal/ws"
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

	fmt.Println("DB connected!")

	if err != nil {
		log.Fatalf("Could not initialize DB table products.")
	}

	// start websocket listener goroutine

	go ws.ListenForWSChannel()

	// Start Server
	app.serve()

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// blocking shutdown until message is sent to the quit channel
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s\n", err)
	}
	log.Println("Server exiting")

}

// shutdown gracefully stops the application
func (app *application) shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		// Close WebSocket connections
		ws.Shutdown()

		// Add other shutdown logic here (e.g., closing DB connections)
		// Simulating some cleanup tasks
		time.Sleep(2 * time.Second)
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown timed out, forcing exit")
		return ctx.Err()
	case <-done:
		log.Println("Shutdown completed successfully")
		return nil
	}
}
