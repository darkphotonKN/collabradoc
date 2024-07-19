package routes

import (
	"net/http"

	"github.com/darkphotonKN/collabradoc/internal/config"
	"github.com/darkphotonKN/collabradoc/internal/docs"
	"github.com/darkphotonKN/collabradoc/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *config.Application) routes() http.Handler {
	mux := chi.NewRouter()

	// handle cors
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*", "http://*",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	/*************************
	* - ROUTES -
	*************************/

	// Docs Routes
	mux.Get("/api/docs", docs.GetDocsList)

	// Users Routes
	mux.Get("/api/user", user.GetUsersHandler)

	return mux
}
