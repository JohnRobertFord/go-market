package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JohnRobertFord/go-market/internal/config"
	"github.com/JohnRobertFord/go-market/internal/handler"
	"github.com/JohnRobertFord/go-market/internal/server"
	"github.com/JohnRobertFord/go-market/internal/storage"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("can't init config: %e", err)
	}
	fmt.Println(cfg)

	ctx := context.Background()
	pg, err := storage.NewStorage(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	userRepository := storage.NewUserRepository(pg)
	userHandler := handler.NewUserHandler(userRepository)

	orderRepository := storage.NewOrderRepository(pg)
	orderHandler := handler.NewOrderHandler(orderRepository)

	go server.NewServer(cfg, userHandler, orderHandler).RunServer()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-sigChan
	log.Println("Shutdown")
}
