package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/JohnRobertFord/go-market/internal/auth"
	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/JohnRobertFord/go-market/internal/storage"
	"github.com/JohnRobertFord/go-market/internal/util"
)

type OrderHandler struct {
	orderRepo *storage.OrderRepository
}

func NewOrderHandler(orderRepo *storage.OrderRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo,
	}
}

// 200 — номер заказа уже был загружен этим пользователем;
// 202 — новый номер заказа принят в обработку;
// 400 — неверный формат запроса;
// 401 — пользователь не аутентифицирован;
// 409 — номер заказа уже был загружен другим пользователем;
// 422 — неверный формат номера заказа;
// 500 — внутренняя ошибка сервера.
func (o *OrderHandler) CreateOrder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	orderID, err := validateOrder(req.Body)
	if err != nil {
		switch err {
		case model.ErrData:
			w.WriteHeader(http.StatusBadRequest)
		case model.ErrCheck:
			w.WriteHeader(http.StatusUnprocessableEntity)
		case model.ErrInternal:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	cookie, err := req.Cookie("Authorization")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := auth.GetUser(cookie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order := &model.Order{
		Username: user,
		Status:   "NEW",
		OrderID:  *orderID,
	}

	err = o.orderRepo.Create(ctx, order)
	if err != nil {
		switch err {
		case model.ErrDataConflict:
			w.WriteHeader(http.StatusConflict)
		case model.ErrDataExist:
			w.WriteHeader(http.StatusOK)
		case model.ErrInternalDB:
			w.WriteHeader(http.StatusInsufficientStorage)
		case model.ErrInternal:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
func (o *OrderHandler) ListOrders(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	cookie, err := req.Cookie("Authorization")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := auth.GetUser(cookie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := o.orderRepo.List(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	out, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func validateOrder(in io.ReadCloser) (*string, error) {

	b, err := io.ReadAll(in)
	if err != nil {
		return nil, model.ErrInternal
	}
	out := string(b)
	bUint, err := strconv.ParseUint(out, 10, 64)
	if err != nil {
		return nil, model.ErrData
	}

	// Luhn check
	if !util.Valid(bUint) {
		return nil, model.ErrCheck
	}

	return &out, nil
}
