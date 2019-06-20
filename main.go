package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Secret struct {
	Hash           string `json:"hash"`
	SecretText     string `json:"secretText"`
	CreatedAt      string `json:"createdAt"`
	ExpiresAt      string `json:"expiresAt"`
	RemainingViews int32  `json:"remainingViews"`
}

var secrets []Secret

func buildSecret(text string, maxViews string, ttl string) Secret {
	remainingViews, err := strconv.ParseInt(maxViews, 10, 32)
	minutes, err := strconv.ParseInt(ttl, 10, 32)

	if err != nil {
		log.Println(err)
	}

	expiresAt := "0"

	if minutes > 0 {
		now := time.Now()
		expiresAt = now.Add(time.Minute * time.Duration(minutes)).String()
	}

	return Secret{
		Hash:           uuid.New().String(),
		SecretText:     text,
		CreatedAt:      time.Now().String(),
		ExpiresAt:      expiresAt,
		RemainingViews: int32(remainingViews),
	}
}

func createSecret(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	request.ParseForm()

	secret := buildSecret(
		request.FormValue("secret"),
		request.FormValue("expireAfterViews"),
		request.FormValue("expireAfter"),
	)
	secrets = append(secrets, secret)

	json.NewEncoder(response).Encode(secret)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/secret", createSecret).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
