package middleware

import (
	"github.com/rs/cors"
	"log"
	"net/http"
)

func Cors(handler http.Handler) http.Handler {
	log.Println("CORS Middleware")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		Debug:            true,
	})
	return c.Handler(handler)
}
