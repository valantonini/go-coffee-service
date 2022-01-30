module github.com/valantonini/go-coffee-service/product-service

go 1.17

replace github.com/valantonini/go-coffee-service/config => ../config

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/matryer/is v1.4.0
	github.com/nats-io/nats.go v1.13.1-0.20211122170419-d7c1d78a50fc
	github.com/valantonini/go-coffee-service/config v1.0.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/nats-io/nats-server/v2 v2.7.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
