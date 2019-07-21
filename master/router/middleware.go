package router

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
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
		log.Println(r.URL)
		user := context.Get(r, "user").(model.User)

		logrus.Info(user.ID.String())

		var count int
		row := model.DB.Raw("SELECT COUNT(*) FROM group_users WHERE "+
			"group_id = (SELECT group_id FROM permissions WHERE endpoint = ? AND verb = ?) "+
			"AND user_id = ?", r.URL.Path, lib.GetVerbByName(r.Method), user.ID.String()).Row()
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
