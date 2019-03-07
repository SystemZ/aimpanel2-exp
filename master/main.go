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
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/middleware"
	"gitlab.com/systemz/aimpanel2/master/router"
	"log"
	"net/http"
)

func main() {
	db.DB = db.SetupDatabase()

	//rabbit.RabbitListen()

	log.Println("Starting API on port :8000")
	r := router.NewRouter()
	log.Fatal(http.ListenAndServe(":8000",
		middleware.CommonMiddleware(middleware.CorsMiddleware(r))))
}
