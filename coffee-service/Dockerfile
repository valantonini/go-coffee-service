FROM golang:1.17-alpine AS builder

WORKDIR /app/

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN go build -o coffee-service .

FROM alpine:3.14 AS runtime

EXPOSE 8080

WORKDIR /app
COPY --from=builder /app/coffee-service /app

ENTRYPOINT ["./coffee-service"]