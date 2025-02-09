package auth

import (
	"log"
	"net/http"
)

var jwtKey = []byte("JWT_Secret")

func AuthJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("AuthJWT" + req.RequestURI)

		h.ServeHTTP(w, req)
	})

}

//	func validateJWT(tokenString string) (*jwt.Token, error) {
//		return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("Unexpected sign method: %v", token.Header["alg"])
//			}
//			return hmacSampleSecret, nil
//		})
//	}
var Claims struct {
	Username string
	jwt.st
}

func CreateJWT() {

}
