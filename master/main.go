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
	db.SetupDatabase()

	r := router.NewRouter()
	log.Fatal(http.ListenAndServe(":9000", middleware.Cors(r)))
	//log.Println("start")
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//lib.FailOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	//
	//channel, err := conn.Channel()
	//lib.FailOnError(err, "Failed to open a channel")
	//defer channel.Close()
	//
	//queue, err := channel.QueueDeclare(
	//	"",
	//	false,
	//	false,
	//	true,
	//	false,
	//	nil)
	//lib.FailOnError(err, "Failed to declare a queue")

	//msgs, err := channel.Consume(
	//	queue.Name,
	//	"",
	//	true,
	//	false,
	//	false,
	//	false,
	//	nil)
	//failOnError(err, "Failed to register a consumer")

	//corrId := lib.RandomString(32)
	//
	//start := lib.RpcMessage{
	//	Type: lib.START,
	//	Body: "alert hello",
	//}
	//jsonMarshal, _ := json.Marshal(start)
	//
	//err = channel.Publish(
	//	"",
	//	"wrapper_rpc",
	//	false,
	//	false,
	//	amqp.Publishing{
	//		ContentType:   "application/json",
	//		CorrelationId: corrId,
	//		ReplyTo:       queue.Name,
	//		Body:          jsonMarshal,
	//	})
	//
	//lib.FailOnError(err, "Failed to publish a message")
}
