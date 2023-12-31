package events

import (
	"github.com/alexandrevicenzi/go-sse"
	"gitlab.com/systemz/aimpanel2/master/config"
	"log"
	"net/http"
	"os"
)

var (
	SSE *sse.Server
)

func InitSSE() *sse.Server {
	options := &sse.Options{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"X-Accel-Buffering":           "no",
		},
	}

	if config.DEV_MODE {
		options.Logger = log.New(os.Stdout, "go-sse: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return sse.NewServer(options)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	SSE.ServeHTTP(w, r)
}
