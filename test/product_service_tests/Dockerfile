FROM golang:1.17 AS builder

WORKDIR /app/
COPY ./go.* ./
RUN go mod download

COPY ./ ./

WORKDIR /app/test/product_service_tests
RUN chmod +x ../wait-for-it.sh ../docker-entrypoint.sh

ENTRYPOINT ["../docker-entrypoint.sh"]

CMD ["go", "test", "--tags=integration", "-v", "./..."]