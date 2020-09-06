package router

import (
	c "context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/exit"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func DBCheckMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := model.DB.Client().Ping(c.TODO(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

		oid, err := primitive.ObjectIDFromHex(token.Claims.(jwt.MapClaims)["uid"].(string))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := model.GetUserById(oid)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		context.Set(r, "user", *user)
		handler.ServeHTTP(w, r)
	})
}

func PermissionMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(model.User)

		route := mux.CurrentRoute(r)
		pathTemplate, _ := route.GetPathTemplate()

		// override all permission system for instance admin
		if user.Admin {
			handler.ServeHTTP(w, r)
			return
		}

		// standard check for standard user
		if model.CheckIfUserHasAccess(r.URL.Path, r.Method, pathTemplate, user.ID) {
			handler.ServeHTTP(w, r)
			return
		}

		logrus.Info("Access denied")
		w.WriteHeader(http.StatusUnauthorized)
		return
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

		hostId, _ := primitive.ObjectIDFromHex(token.Claims.(jwt.MapClaims)["uid"].(string))
		host, err := model.GetHostById(hostId)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if host == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := mux.Vars(r)
		if params["hostToken"] != host.Token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		gsId, ok := params["gsId"]
		if ok {
			oid, err := primitive.ObjectIDFromHex(gsId)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			gs, err := model.GetGameServerById(oid)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if gs == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		context.Set(r, "host", *host)
		handler.ServeHTTP(w, r)
	})
}

func CorsMiddleware(handler http.Handler) http.Handler {
	log.Println("CORS Middleware")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Accept", "Content-Length", "X-Requested-With", "Origin"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		Debug:          false,
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
