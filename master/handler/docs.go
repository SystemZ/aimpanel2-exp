package handler

import (
	"net/http"
)

func GetSpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger.json")
}

func GetDocsRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/v1/docs/", 302)
}

func GetDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "redoc.html")
}

func GetSwaggerUi(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger.html")
}
