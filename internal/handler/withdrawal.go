package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/JohnRobertFord/go-market/internal/auth"
	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/JohnRobertFord/go-market/internal/storage"
)

type WithdrawalHandler struct {
	withdrawalRepo *storage.WithdrawalRepository
}

func NewWithdrawalHandler(withdrawalRepo *storage.WithdrawalRepository) *WithdrawalHandler {
	return &WithdrawalHandler{
		withdrawalRepo,
	}
}

// 200 — успешная обработка запроса;
// 401 — пользователь не авторизован;
// 402 — на счету недостаточно средств;
// 422 — неверный номер заказа;
// 500 — внутренняя ошибка сервера.
func (wh *WithdrawalHandler) RequestWithdraw(w http.ResponseWriter, req *http.Request) {
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

	data, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var in *model.WithdrawalRequest
	err = json.Unmarshal(data, in)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = wh.withdrawalRepo.Withdraw(ctx, user, in)
	if err != nil {
		switch err {
		case model.ErrInsufficientFunds:
			w.WriteHeader(http.StatusPaymentRequired)
		case model.ErrCheck:
			w.WriteHeader(http.StatusUnprocessableEntity)
		case model.ErrInternal:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// 200 — успешная обработка запроса.
// 204 — нет ни одного списания.
// 401 — пользователь не авторизован.
// 500 — внутренняя ошибка сервера.
func (wh *WithdrawalHandler) GetWithdrawals(w http.ResponseWriter, req *http.Request) {
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

	res, err := wh.withdrawalRepo.Withdrawals(ctx, user)
	if err != nil {
		if err == model.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
