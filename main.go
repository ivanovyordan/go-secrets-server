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
		end, _ := now.Add(time.Minute * time.Duration(minutes)).MarshalText()
		expiresAt = string(end)
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

func getSecret(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)

	for index, secret := range secrets {
		if secret.Hash == params["hash"] && secret.RemainingViews > 0 {
			parsed, _ := time.Parse(time.RFC3339, secret.ExpiresAt)

			if time.Now().After(parsed) {
				break
			}

			secret.RemainingViews -= 1
			secrets[index].RemainingViews -= 1
			json.NewEncoder(response).Encode(secret)
			return
		}
	}

	json.NewEncoder(response).Encode(&Secret{})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/secret", createSecret).Methods("POST")
	router.HandleFunc("/v1/secret/{hash}", getSecret).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
