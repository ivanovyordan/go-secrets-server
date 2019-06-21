package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var all []Secret

type Secret struct {
	Hash           string `json:"hash" xml:"hash"`
	SecretText     string `json:"secretText" xml:"secretText"`
	CreatedAt      string `json:"createdAt" xml:"createdAt"`
	ExpiresAt      string `json:"expiresAt" xml:"expiresAt"`
	RemainingViews int32  `json:"remainingViews" xml:"remainingViews"`
}

func NewSecret(text string, maxViews string, ttl string) (Secret, error) {
	remainingViews, err := strconv.ParseInt(maxViews, 10, 32)
	if err != nil {
		return invalidInput()
	}

	minutes, err := strconv.ParseInt(ttl, 10, 32)
	if err != nil {
		return invalidInput()
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

	return secret, nil
}

func FindSecret(hash string) (Secret, error) {
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

		return secret, nil
	}

	return Secret{}, errors.New("Secret not found")
}

func invalidInput() (Secret, error) {
	return Secret{}, errors.New("Invalid input")
}
