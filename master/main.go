//go:generate swagger generate spec

// Package classification Aimpanel Master API
//
// Schemes: http, https
// Host: localhost:9000
// BasePath: /v1
// Version: 0.0.1
//
// Consumes:
// 	- application/json
//
// Produces:
// 	- application/json
//
// swagger:meta
package main

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/middleware"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
	"gitlab.com/systemz/aimpanel2/master/redis"
	"gitlab.com/systemz/aimpanel2/master/router"
	"net/http"
)

func main() {
	db.Setup()
	redis.Setup()

	rabbit.Listen()
	rabbit.ListenWrapperLogsQueue()

	logrus.Info("Starting API on port :8000")
	r := router.NewRouter()
	logrus.Fatal(http.ListenAndServe(":8000",
		middleware.CommonMiddleware(middleware.CorsMiddleware(r))))
}
