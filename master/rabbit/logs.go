package rabbit

import (
	"gitlab.com/systemz/aimpanel2/lib"
	"log"
)

func ListenWrapperLogsQueue() {
	msgs, err := channel.Consume(
		"wrapper_logs",
		"",
		true,
		false,
		false,
		false,
		nil)
	lib.FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}
