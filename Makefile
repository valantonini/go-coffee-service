test:
	cd coffee-service && go mod download
	cd coffee-service && go test -v ./...

integration:
	docker compose -f docker-compose.yml -f docker-compose.integration.yml up --build --abort-on-container-exit --remove-orphans
