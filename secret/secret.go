package secret

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var all []Secret

type Secret struct {
	Hash           string `json:"hash"`
	SecretText     string `json:"secretText"`
	CreatedAt      string `json:"createdAt"`
	ExpiresAt      string `json:"expiresAt"`
	RemainingViews int32  `json:"remainingViews"`
}

func Build(text string, maxViews string, ttl string) Secret {
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

	secret := Secret{
		Hash:           uuid.New().String(),
		SecretText:     text,
		CreatedAt:      time.Now().String(),
		ExpiresAt:      expiresAt,
		RemainingViews: int32(remainingViews),
	}

	all = append(all, secret)

	return secret
}

func Find(hash string) Secret {
	for index, secret := range all {
		if secret.Hash != hash {
			continue
		}

		expirationTime, _ := time.Parse(time.RFC3339, secret.ExpiresAt)
		if secret.RemainingViews == 0 || time.Now().After(expirationTime) {
			break
		}

		secret.RemainingViews -= 1
		all[index] = secret

		return secret
	}

	return Secret{}
}
