FROM golang:1.17-alpine AS builder

WORKDIR /config
COPY ./config/ ./

WORKDIR /app/
COPY product-service/go.* ./
RUN go mod download
COPY ./product-service/ ./

RUN go build -o product-service .

FROM alpine:3.14 AS runtime

EXPOSE 9000

WORKDIR /app
COPY --from=builder /app/product-service /app

ENTRYPOINT ["./product-service"]