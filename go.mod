module github.com/valantonini/go-coffee-service

go 1.17

replace github.com/valantonini/go-coffee-service/config => ./internal/pkg/config

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/matryer/is v1.4.0
	github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d
	github.com/valantonini/go-coffee-service/order-service v0.0.0-20220210203804-f78582ef3d77
	github.com/valantonini/go-coffee-service/product-service v0.0.0-20220210203804-f78582ef3d77
)

require (
	github.com/klauspost/compress v1.14.2 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/nats-server/v2 v2.7.2 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
)
