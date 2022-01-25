MODULE_DIRS =  ./product-service ./order-service ./config

restore:
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go mod download) &&) true

test:
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go test -v ./...) &&) true

integration:
	docker compose -f docker-compose.yml -f docker-compose.integration.yml up --build --abort-on-container-exit --remove-orphans

test_all: test integration