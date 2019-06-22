package main

import (
	"log"
	"net/http"

	"secrets/controller/secret"
	"secrets/tools/db"
	"secrets/tools/metrics"

	"github.com/gorilla/mux"
)

func init() {
	metrics.Init()
}

func main() {
	router := mux.NewRouter()
	router.Use(metrics.Middleware)
	err := db.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

	router.Handle("/metrics", metrics.Get()).Methods("GET").Name("GetMetrics")
	router.HandleFunc("/secret", secret.Post).Methods("POST").Name("PostSecret")
	router.HandleFunc("/secret/{hash}", secret.Get).Methods("GET").Name("GetSecret")

	log.Fatal(http.ListenAndServe(":8000", router))
}
