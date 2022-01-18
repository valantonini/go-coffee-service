# go-coffee-service

Playing with patterns and libraries for distributed systems in the Go ecosystem.
Based on the [Hashicorp consul tutorial](https://learn.hashicorp.com/tutorials/consul/kubernetes-extract-microservice?in=consul/microservices)

## caveats

- i am new to go
- striving to use as few dependencies as possible to understand the problems libraries and frameworks are addressing 
- secret management will not be production ready


## env

- go 1.17
- docker 20.10.11
- make

```bash
docker network create -d bridge coffee-service-network

```

## run

```bash
# run the integration tests
make integration

# run the coffee-service
docker compose up

```
