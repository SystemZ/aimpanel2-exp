package sse

import (
	"github.com/alexandrevicenzi/go-sse"
	"net/http"
)

var (
	SSE *sse.Server
)


func InitSSE() *sse.Server {
	return sse.NewServer(&sse.Options{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	})
}

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	SSE.ServeHTTP(w, r)
}