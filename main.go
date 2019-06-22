package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"secrets/model"
	"secrets/tools"

	"github.com/gorilla/mux"
)

const (
	XMLHeader = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

func respond(response http.ResponseWriter, request *http.Request, message model.Secret) {
	accept := strings.ToLower(request.Header.Get("Accept"))
	var content []byte

	if accept == "application/xml" {
		response.Header().Set("Content-Type", "application/xml")
		content, _ = xml.MarshalIndent(message, "", "  ")
		content = []byte(XMLHeader + string(content))

	} else if accept == "application/json" {
		response.Header().Set("Content-Type", "application/json")
		content, _ = json.Marshal(message)
	} else {
		fail(response, http.StatusNotAcceptable, "")
	}

	response.Write(content)
}

func fail(response http.ResponseWriter, code int, message string) {
	http.Error(response, message, code)
}

func postSecret(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	message, err := model.NewSecret(
		request.FormValue("secret"),
		request.FormValue("expireAfterViews"),
		request.FormValue("expireAfter"),
	)

	if err != nil {
		fail(response, http.StatusMethodNotAllowed, err.Error())
		return
	}

	respond(response, request, message)
}

func getSecret(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	message, err := model.FindSecret(params["hash"])

	if err != nil {
		fail(response, http.StatusNotFound, err.Error())
		return
	}

	respond(response, request, message)
}

func init() {
	tools.InitMetrics()
}

func main() {
	router := mux.NewRouter()
	router.Use(tools.MetricsMiddleware)

	router.Handle("/metrics", tools.GetMetrics()).Methods("GET").Name("GetMetrics")
	router.HandleFunc("/v1/secret", postSecret).Methods("POST").Name("PostSecret")
	router.HandleFunc("/v1/secret/{hash}", getSecret).Methods("GET").Name("GetSecret")

	log.Fatal(http.ListenAndServe(":8000", router))
}
