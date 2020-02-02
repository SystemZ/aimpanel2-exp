package router

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/exit"
	"gitlab.com/systemz/aimpanel2/master/model"
	"log"
	"net/http"
	"os"
	"strings"
)

func CommonMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}

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

		var user model.User
		if model.DB.Where("id = ?", token.Claims.(jwt.MapClaims)["uid"].(string)).First(&user).RecordNotFound() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		context.Set(r, "user", user)

		handler.ServeHTTP(w, r)
	})
}

func PermissionMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(model.User)

		var count int
		row := model.DB.Raw("SELECT COUNT(*) FROM permissions WHERE "+
			"endpoint = ? AND verb = ? AND group_id = (SELECT group_id FROM group_users WHERE user_id = ?);",
			r.URL.Path, lib.GetVerbByName(r.Method), user.ID.String()).Row()
		err := row.Scan(&count)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if count == 0 {
			logrus.Info("Access denied")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func SlavePermissionMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//TODO: move to config package instead of jwt_secret
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var host model.Host
		if model.DB.Where("id = ?", token.Claims.(jwt.MapClaims)["uid"].(string)).First(&host).RecordNotFound() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		if params["host_token"] != host.Token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		gsId, ok := params["server_id"]
		if ok {
			var gs model.GameServer
			if model.DB.Where("id = ?", gsId).First(&gs).RecordNotFound() {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

func CorsMiddleware(handler http.Handler) http.Handler {
	log.Println("CORS Middleware")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Accept", "Content-Length", "X-Requested-With", "Origin"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		Debug:          true,
	})
	return c.Handler(handler)
}

func ExitMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if exit.EXIT {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
