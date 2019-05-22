//go:generate swagger generate spec

// Package classification Aimpanel Master API
//
// Schemes: http, https
// Host: localhost:8000
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
	"gitlab.com/systemz/aimpanel2/master/cmd"
)

func main() {
	cmd.Execute()
}
