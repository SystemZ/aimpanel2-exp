package handler

import (
	"gitlab.com/systemz/aimpanel2/master/config"
	"net/http"
)

func GetDocsRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/v1/docs/", 302)
}

func GetSpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.HTTP_DOCS_DIR+"swagger.json")
}

func GetDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, config.HTTP_DOCS_DIR+"redoc.html")
}

func GetSwaggerUi(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.HTTP_DOCS_DIR+"swagger.html")
}
