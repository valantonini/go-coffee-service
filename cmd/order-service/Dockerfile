FROM golang:1.17-alpine AS builder

WORKDIR /app/
COPY ./go.* ./
RUN go mod download
COPY ./ ./

WORKDIR /app/cmd/order-service
RUN go build -o order-service .

FROM alpine:3.14 AS runtime

EXPOSE 9001

WORKDIR /app
COPY --from=builder /app/cmd/order-service /app

ENTRYPOINT ["./order-service"]