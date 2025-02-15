package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JohnRobertFord/go-market/internal/accrual"
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

	balanceRepository := storage.NewBalanceRepository(pg)
	balanceHandler := handler.NewBalanceHandler(balanceRepository)

	withdrawalRepository := storage.NewWithdrawalRepository(pg)
	withdrawalHandler := handler.NewWithdrawalHandler(withdrawalRepository)

	accrualRepository := storage.NewAccuralRepository(pg)
	accrualService := accrual.NewWorker(cfg.Accrual, accrualRepository)

	go server.NewServer(
		cfg,
		userHandler,
		orderHandler,
		balanceHandler,
		withdrawalHandler,
	).RunServer()

	go accrualService.Run(ctx)

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
