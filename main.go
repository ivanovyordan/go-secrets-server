package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"secrets/secret"
)

func respond(response http.ResponseWriter, request *http.Request, message secret.Secret) {
	accept := strings.ToLower(request.Header.Get("Accept"))
	var content []byte

	if accept == "application/xml" {
		const (
			Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
		)

		response.Header().Set("Content-Type", "application/xml")
		content, _ = xml.MarshalIndent(message, "", "  ")
		content = []byte(Header + string(content))

	} else if accept == "application/json" {
		response.Header().Set("Content-Type", "application/json")
		content, _ = json.Marshal(message)
	}

	response.Write(content)
}

func createSecret(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	message := secret.Build(
		request.FormValue("secret"),
		request.FormValue("expireAfterViews"),
		request.FormValue("expireAfter"),
	)

	respond(response, request, message)
}

func getSecret(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	message := secret.Find(params["hash"])

	respond(response, request, message)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/secret", createSecret).Methods("POST")
	router.HandleFunc("/v1/secret/{hash}", getSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
