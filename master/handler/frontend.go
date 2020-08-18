package handler

import (
	"gitlab.com/systemz/aimpanel2/master/config"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, config.HTTP_FRONTEND_DIR+"index.html")
}

func ServiceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	//w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	http.ServeFile(w, r, config.HTTP_FRONTEND_DIR+"service-worker.js")
}

func RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	http.ServeFile(w, r, config.HTTP_FRONTEND_DIR+"robots.txt")
}
