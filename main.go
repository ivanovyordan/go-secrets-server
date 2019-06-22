package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ivanovyordan/go-secrets-server/controller/secret"
	"github.com/ivanovyordan/go-secrets-server/tools/db"
	"github.com/ivanovyordan/go-secrets-server/tools/metrics"

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

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
