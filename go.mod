module github.com/valantonini/go-coffee-service

go 1.17

replace github.com/valantonini/go-coffee-service/cmd/product-service => ./cmd/product-service

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v1.1.0
	github.com/matryer/is v1.4.0
	github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d
	go.mongodb.org/mongo-driver v1.8.3
)

require (
	github.com/fatih/color v1.7.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.14.2 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/nats-server/v2 v2.7.2 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
)
