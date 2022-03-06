# go-coffee-service

[![CI](https://github.com/valantonini/go-coffee-service/actions/workflows/makefile.yml/badge.svg)](https://github.com/valantonini/go-coffee-service/actions/workflows/makefile.yml)

playing with patterns and libraries for distributed systems in the Go ecosystem.

inspired by the [Hashicorp consul tutorial](https://learn.hashicorp.com/tutorials/consul/kubernetes-extract-microservice?in=consul/microservices) and [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

## progress

### patterns
outbox: guarantees _at least once delivery_ of messaging by first writing events to an outbox table within transactions that mutate the database. A Go routine periodically polls the outbox
table, dispatches unsent messages onto the nats message bus and updates the outbox entries as sent.

### product-service 
a webapi in charge of coffee listings (GetAll/GetById/Add). service implements the outbox pattern and raises events on nats broker after db mutations (e.g. adding a new coffee)

### order-service
a webapi that queries for a list of coffees from the product-service on startup via message bus (nats) 

## caveats

- i am new to go and many things aren't in their final shape
- striving to use as few dependencies as possible to understand the problems libraries and frameworks are addressing 
- secret management will not be production ready
- to do:
  - extract inline variables out to env
  - durable messaging
  - message retry

## env

- go 1.17
- docker 20.10.11
- make

```bash
docker network create -d bridge coffee-service-network
```

## run

```bash
# run the coffee-services
docker compose up

# run unit and integration tests
make test_all
```


## optional tooling for debugging

nats realtime monitoring

```bash
go get github.com/nats-io/nats-top

nats-top
```

nats command line tool (osx)
```bash
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats
```
