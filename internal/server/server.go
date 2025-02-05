package server

import (
	"context"
	"log"
	"net/http"

	"github.com/JohnRobertFord/go-market/internal/config"
	"github.com/JohnRobertFord/go-market/internal/handler"
	"github.com/JohnRobertFord/go-market/internal/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	// "github.com/go-chi/chi/v5/middleware"
)

type server struct {
	Server *http.Server
}

func (s server) RunServer() {
	log.Fatal(s.Server.ListenAndServe())
	s.Server.Shutdown(context.Background())
}

func NewServer(cfg *config.Config, userHandler *handler.UserHandler) *server {

	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"), logger.Logging)

	r.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/orders", handler.Placeholder())
		r.Get("/orders", handler.Placeholder())
		r.Get("/balance", handler.Placeholder())
		r.Post("/balance/withdraw", handler.Placeholder())
		r.Get("/withdrawals", handler.Placeholder())
	})

	return &server{
		Server: &http.Server{
			Addr:    cfg.Bind,
			Handler: r,
		},
	}
}
