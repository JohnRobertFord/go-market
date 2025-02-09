package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/JohnRobertFord/go-market/internal/storage"
)

type UserHandler struct {
	userRepo *storage.UserRepository
}
type registerRequest struct {
	Name string `json:"login" `
	PWD  string `json:"password"`
}

var (
	succes = map[string]string{
		"message": "Success",
	}
)

func NewUserHandler(userRepo *storage.UserRepository) *UserHandler {
	return &UserHandler{userRepo}
}

func Placeholder() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// ctx := req.Context()

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Placeholder\n")

	})
}

func (uh *UserHandler) Register(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	decoder := json.NewDecoder(req.Body)
	var in registerRequest
	err := decoder.Decode(&in)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	defer req.Body.Close()

	if len(in.Name) == 0 || len(in.PWD) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = uh.userRepo.NewUser(ctx, &model.User{Name: in.Name, Password: in.PWD})
	if err != nil {
		if err == model.ErrDataConflict {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) Login(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	decoder := json.NewDecoder(req.Body)
	var in registerRequest
	err := decoder.Decode(&in)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	defer req.Body.Close()

	if len(in.Name) == 0 || len(in.PWD) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = uh.userRepo.ValidateUser(ctx, &model.User{Name: in.Name, Password: in.PWD})
	if err != nil {
		if err == model.ErrValidate {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func respondeJSON(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&message); err != nil {
		fmt.Println(err)
		return
	}
}
