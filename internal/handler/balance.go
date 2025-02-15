package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JohnRobertFord/go-market/internal/auth"
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
	user, err := auth.GetUser(req)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := b.balanceRepo.Balance(ctx, user)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
