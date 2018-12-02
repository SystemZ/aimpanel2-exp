package handler

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET / index Index
	//
	// Just index page.
	//
	// This will show nothing.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200:

	w.WriteHeader(http.StatusOK)
}
