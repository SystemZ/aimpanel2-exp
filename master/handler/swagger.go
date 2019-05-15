package handler

import (
	"io/ioutil"
	"log"
	"net/http"
)

func SwaggerSpec(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("swagger.json")
	if err != nil {
		log.Printf("%v", err)
	}
	w.Write(body)
}
