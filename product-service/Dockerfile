FROM golang:1.17-alpine AS builder

WORKDIR /app/

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN go build -o product-service .

FROM alpine:3.14 AS runtime

EXPOSE 9000

WORKDIR /app
COPY --from=builder /app/product-service /app

ENTRYPOINT ["./product-service"]