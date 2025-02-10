package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/JohnRobertFord/go-market/internal/auth"
	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/JohnRobertFord/go-market/internal/storage"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("JWT_Secret")

type (
	UserHandler struct {
		userRepo *storage.UserRepository
	}
	registerRequest struct {
		Name string `json:"login" `
		PWD  string `json:"password"`
	}
	Claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}
)

func NewUserHandler(userRepo *storage.UserRepository) *UserHandler {
	return &UserHandler{userRepo}
}

func Placeholder() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

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

	cookie, err := auth.CreateJWT(in.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)
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
		switch err {
		case model.ErrValidate:
			w.WriteHeader(http.StatusUnauthorized)
		case model.ErrInternal:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	cookie, err := auth.CreateJWT(in.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
