package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set("uid", token.Claims.(jwt.MapClaims)["uid"].(string))

		handler.ServeHTTP(w, r)
	})
}
