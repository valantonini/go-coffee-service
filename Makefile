MODULE_DIRS =  ./cmd/product-service ./cmd/order-service ./config

restore:
	go mod download

test:
	go test -v ./...

integration:
	docker compose -f docker-compose.yml -f docker-compose.integration.yml up --build --abort-on-container-exit --remove-orphans

test_all: test integration