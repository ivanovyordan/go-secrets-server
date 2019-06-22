package secret

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"secrets/tools/db"
	"strconv"
	"time"
)

type Secret struct {
	Hash           string `json:"hash" xml:"hash"`
	SecretText     string `json:"secretText" xml:"secretText"`
	CreatedAt      string `json:"createdAt" xml:"createdAt"`
	ExpiresAt      string `json:"expiresAt" xml:"expiresAt"`
	RemainingViews int32  `json:"remainingViews" xml:"remainingViews"`
}

func New(text string, maxViews string, ttl string) (Secret, error) {
	if !isValid(text, maxViews, ttl) {
		return Secret{}, errors.New("Invalid input")
	}

	now := time.Now()
	nowText, _ := now.MarshalText()
	expiresAt := expirationTime(ttl)
	remainingViews, _ := strconv.ParseInt(maxViews, 10, 32)
	hash := buildHash(text, maxViews, string(nowText), ttl)

	_, err := db.Connection.Exec(
		"INSERT INTO secrets (hash, secret_text, expires_at, remaining_views) VALUES ($1, $2, $3, $4)",
		hash, text, expiresAt, remainingViews,
	)

	if err != nil {
		return Secret{}, err
	}

	return find(hash)
}

func Find(hash string) (Secret, error) {
	secret, err := find(hash)

	if err == nil {
		decreaseViews(&secret)
	}

	return secret, err
}

func decreaseViews(secret *Secret) {
	secret.RemainingViews -= 1
	db.Connection.Exec(`UPDATE secrets SET remaining_views = remaining_views - 1 WHERE hash = $1`, secret.Hash)
}

func find(hash string) (Secret, error) {
	var hashText string
	var secretText string
	var createdAt string
	var expiresAt string
	var remainingViews int32

	err := db.Connection.QueryRow(`
		SELECT
			hash,
			secret_text,
			created_at,
			remaining_views,
			CASE
				WHEN expires_at = '0001-01-01 00:00:00' THEN '0'
				ELSE to_char(expires_at, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')
			END AS expires_at
		FROM secrets
		WHERE
			hash = $1
			AND remaining_views > 0
			AND (expires_at = '0001-01-01 00:00:00' OR expires_at > CURRENT_TIMESTAMP)
		LIMIT 1
	`, hash).Scan(&hashText, &secretText, &createdAt, &remainingViews, &expiresAt)

	if err != nil {
		return Secret{}, errors.New("Secret not found")
	}

	secret := Secret{
		Hash:           hashText,
		SecretText:     secretText,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
		RemainingViews: remainingViews,
	}

	return secret, nil
}

func isValid(text string, maxViews string, ttl string) bool {
	if text == "" || maxViews == "" || ttl == "" {
		return false
	}

	remainingViews, err := strconv.ParseInt(maxViews, 10, 32)
	if err != nil || remainingViews < 1 {
		return false
	}

	minutes, err := strconv.ParseInt(ttl, 10, 32)
	if err != nil || minutes < 0 {
		return false
	}

	return true
}

func expirationTime(ttl string) time.Time {
	var expiresAt time.Time
	minutes, _ := strconv.ParseInt(ttl, 10, 32)

	if minutes > 0 {
		expiresAt = time.Now().UTC().Add(time.Minute * time.Duration(minutes))
	}

	return expiresAt
}

func buildHash(text string, maxViews string, createdAt string, ttl string) string {
	data := text + maxViews + createdAt + ttl
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
