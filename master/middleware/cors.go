package middleware

import (
	"github.com/rs/cors"
	"log"
	"net/http"
)

func CorsMiddleware(handler http.Handler) http.Handler {
	log.Println("CORS Middleware")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "Content-Length"},
		Debug:            true,
	})
	return c.Handler(handler)
}
