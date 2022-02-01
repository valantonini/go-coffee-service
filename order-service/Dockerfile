FROM golang:1.17-alpine AS builder

# copy config over to ../config (relative to app) for local module replacements
WORKDIR /config
COPY ./config/ ./

WORKDIR /app/
COPY order-service/go.* ./
RUN go mod download
COPY ./order-service/ ./

RUN go build -o order-service .

FROM alpine:3.14 AS runtime

EXPOSE 9001

WORKDIR /app
COPY --from=builder /app/order-service /app

ENTRYPOINT ["./order-service"]