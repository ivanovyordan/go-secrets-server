module github.com/ivanovyordan/go-secrets-server

go 1.12

require (
	github.com/gorilla/mux v1.7.3
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.1.1
	github.com/prometheus/client_golang v1.1.0
)

replace github.com/ivanovyordan/go-secrets-server => ./
