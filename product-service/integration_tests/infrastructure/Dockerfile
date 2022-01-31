FROM golang:1.17 AS builder

# copy config over to ../config (relative to app) for local module replacements
WORKDIR /config
COPY ./config/ ./

WORKDIR /app/
COPY product-service/go.* ./
RUN go mod download

COPY ./product-service/ ./

RUN chmod +x ./integration_tests/infrastructure/wait-for-it.sh ./integration_tests/infrastructure/docker-entrypoint.sh

ENTRYPOINT ["/app/integration_tests/infrastructure/docker-entrypoint.sh"]

WORKDIR /app/integration_tests

CMD ["go", "test", "--tags=integration", "-v", "./..."]