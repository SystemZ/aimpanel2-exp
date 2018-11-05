package middleware

import (
	"github.com/rs/cors"
	"log"
	"net/http"
)

func Cors(handler http.Handler) http.Handler {
	log.Println("CORS Middleware")
	return cors.Default().Handler(handler)
}
