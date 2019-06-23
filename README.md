# Go Secrets Server
A simple http API for storing and accessing secrets.

## Features
The secret server can be used to store and share secrets using the random generated URL.
But the secret can be read only a limited number of times after that it will expire and won’t be available.
The secret may have TTL.
After the expiration time the secret won't be available anymore.

## Requirements
* PostgreSQL server

## Installation
1. Clone the repository
2. Install the dependencies `go mod vendor`
3. Copy the `.env.example` file to `.env` and change the variables there
4. Create the `secrets` table using the SQL from `create_table.sql`

## Usage
Root URL: `https://secrets-ivanovyordan.herokuapp.com/`.
The server has 3 simple API endpoints.

### Metrics
* **Endpoint:** `/metrics`
* **Method:** `GET`
* **Description:** A Prometheus endpoint with the following custom metrics:
	* `http_request_endpoint_calls_total` - Request counter for each API endpoint
	* `http_request_endpoint_duration_seconds` - Response time for each API endpoint by buckets
* **Example:** `http https://secrets-ivanovyordan.herokupp.com/metrics`

### Create a secret
* **Endpoint:** `/secret`
* **Method:** `POST`
* **Headers:** The endpoint expects and `ACCEPT` header with either `application/json` or `application/xml` value. The request will fail if no `ACCEPT` header is passed.
* **Parameters:** The endpoint expects all of the following parameters passed as form-data:
	* `secret` - This text will be saved as a secret
	* `expireAfterViews` - The secret won't be available after the given number of views. It must be greater than 0
	* `expireAfter` - The secret won't be available after the given time. The value is provided in minutes. 0 means never expires
* **Example:** `http --form POST https://secrets-ivanovyordan.herokupp.com/secret Accept:application/json secret="My little secret" expireA fterViews=10 expireAfter=10`

### Get a secret
* **Endpoint:** `/secret/<hash>`
* **Method:** `GET`
* **Headers:** The endpoint expects and `ACCEPT` header with either `application/json` or `application/xml` value. The request will fail if no `ACCEPT` header is passed.
* **Parameters:**
	* `hash` - Unique hash to identify the secret
* **Example:** `http https://secrets-ivanovyordan.herokupp.com/secret/68c418c4066204460e6c1139f37e855e2d4c5b77849b567625a0c7dfe8324ced Accept:application/json`
