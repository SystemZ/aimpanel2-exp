package handler

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
