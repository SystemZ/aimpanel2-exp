package handler

import (
	"gitlab.com/systemz/aimpanel2/master/config"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.HTTP_FRONTEND_DIR+"index.html")
}
