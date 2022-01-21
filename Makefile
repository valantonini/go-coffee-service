MODULE_DIRS =  ./coffee-service

test:
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go test -v ./...) &&) true

integration:
	docker compose -f docker-compose.yml -f docker-compose.integration.yml up --build --abort-on-container-exit --remove-orphans
