package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

func Cors(handler http.Handler) http.Handler {
	return cors.Default().Handler(handler)
}
