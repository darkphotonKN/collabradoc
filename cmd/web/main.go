package main

import (
	"fmt"
	"net/http"
	"time"
)

type config struct {
	port int
}

type application struct {
	config config
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
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
	config := config{
		port: 5000,
	}

	app := &application{
		config: config,
	}

	// start server
	app.serve()
}
