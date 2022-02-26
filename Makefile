restore:
	go mod download

unit:
	go test -v ./...

integration:
	docker compose -f docker-compose.yml -f docker-compose.integration.yml up --build --abort-on-container-exit --remove-orphans

test_all: unit integration