package server

import (
	"context"
	"log"
	"net/http"

	"github.com/JohnRobertFord/go-market/internal/auth"
	"github.com/JohnRobertFord/go-market/internal/config"
	"github.com/JohnRobertFord/go-market/internal/handler"
	"github.com/JohnRobertFord/go-market/internal/logger"
	"github.com/go-chi/chi"
)

type server struct {
	Server *http.Server
}

func (s server) RunServer() {
	log.Fatal(s.Server.ListenAndServe())
	s.Server.Shutdown(context.Background())
}

func NewServer(cfg *config.Config, userHandler *handler.UserHandler, orderHandler *handler.OrderHandler, balanceHandler *handler.BalanceHandler) *server {

	r := chi.NewRouter()
	r.Use(logger.Logging, auth.AuthJWT)

	r.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Post("/orders", orderHandler.CreateOrder)
		r.Get("/orders", orderHandler.ListOrders)
		r.Get("/balance", balanceHandler.GetBalance)
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
