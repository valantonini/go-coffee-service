FROM golang:1.17 AS builder

WORKDIR /app/

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN chmod +x ./integration_tests/infrastructure/wait-for-it.sh ./integration_tests/infrastructure/docker-entrypoint.sh $COFFEE_SERVICE_URL

ENV INTEGRATION_TESTS=1

ENTRYPOINT ["/app/integration_tests/infrastructure/docker-entrypoint.sh"]

WORKDIR /app/integration_tests

CMD ["go", "test", "-v", "./..."]