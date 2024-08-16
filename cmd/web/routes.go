package main

import (
	"net/http"
	"time"

	"github.com/darkphotonKN/collabradoc/internal/comment"
	"github.com/darkphotonKN/collabradoc/internal/document"
	livesession "github.com/darkphotonKN/collabradoc/internal/live_session"
	"github.com/darkphotonKN/collabradoc/internal/user"
	"github.com/darkphotonKN/collabradoc/internal/utils/auth"
	"github.com/darkphotonKN/collabradoc/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
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

	// handles and recovers from any Panic and responds with the error code
	mux.Use(middleware.Recoverer)

	// prevent any command or request taking more than 60 seconds to time out
	mux.Use(middleware.Timeout(60 * time.Second))

	/*************************
	* - ROUTES -
	*************************/

	// -- WebSocket Routes --
	mux.Get("/ws", ws.WsHandler)

	// -- Docs Routes --
	mux.With(auth.JWTMiddleware).Get("/api/doc", document.GetDocumentsHandler)
	mux.With(auth.JWTMiddleware).Post("/api/doc", document.CreateDocHandler)

	// -- Live Session Routes --
	mux.With(auth.JWTMiddleware).Post("/api/livesession", livesession.CreateLiveSessionHandler)
	mux.With(auth.JWTMiddleware).Get("/api/livesession", livesession.GetLiveSessionLinkHandler)

	// -- Comment Routes --
	mux.With(auth.JWTMiddleware).Post("/api/comment", comment.CreateCommentHandler)

	// -- Users Routes --
	mux.With(auth.JWTMiddleware).Get("/api/user", user.GetUsersHandler)
	mux.Post("/api/user/signup", user.SignUpHandler)
	mux.Post("/api/user/login", user.LoginHandler)

	return mux
}
