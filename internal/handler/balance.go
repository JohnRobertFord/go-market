package handler

import (
	"io"
	"net/http"

	"github.com/JohnRobertFord/go-market/internal/storage"
)

type BalanceHandler struct {
	balanceRepo *storage.BalanceRepository
}

func NewBalanceHandler(balanceRepo *storage.BalanceRepository) *BalanceHandler {
	return &BalanceHandler{
		balanceRepo,
	}
}

func (b *BalanceHandler) GetBalance(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	b.balanceRepo.Balance(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Placeholder\n")
}
