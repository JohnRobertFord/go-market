package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("JWT_Secret")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.RequestURI, "login") && !strings.Contains(req.RequestURI, "register") {
			cookie, err := req.Cookie("Authorization")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tokenStr := cookie.Value
			claims := &Claims{}
			tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				if _, ok := err.(*jwt.ValidationError); ok {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
		h.ServeHTTP(w, req)
	})

}

func CreateJWT(username string) (*http.Cookie, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, model.ErrInternal
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: expirationTime,
	}

	return cookie, nil
}

func GetUser(cookie *http.Cookie) (string, error) {
	tokenStr := cookie.Value
	claims := &Claims{}
	jwt.ParseWithClaims(tokenStr, claims, nil)

	return claims.Username, nil
}
