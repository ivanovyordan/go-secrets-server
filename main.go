package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"secrets/secret"
)

func createSecret(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	request.ParseForm()

	message := secret.Build(
		request.FormValue("secret"),
		request.FormValue("expireAfterViews"),
		request.FormValue("expireAfter"),
	)

	json.NewEncoder(response).Encode(message)
}

func getSecret(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)
	message := secret.Find(params["hash"])

	json.NewEncoder(response).Encode(message)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/secret", createSecret).Methods("POST")
	router.HandleFunc("/v1/secret/{hash}", getSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
