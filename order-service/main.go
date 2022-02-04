package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/config"
	"github.com/valantonini/go-coffee-service/order-service/gateway"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfigFromEnv()

	cfg.Logger.Printf("connecting to nats on %v\n", cfg.NatsAddress)
	// var b gateway.Bus
	nc, err := nats.Connect(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Fatal(err)
	}
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		cfg.Logger.Fatal(err)
	}
	cfg.Logger.Println("connected to nats")

	cfg.Logger.Println("retrieving coffee from coffee-service")
	coffeeService := gateway.NewCoffeeServiceGateway(ec)
	coffees, err := coffeeService.GetAll()
	if err != nil {
		cfg.Logger.Println(err)
	}
	cfg.Logger.Printf("%#v\n", coffees)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(fmt.Sprintf("%#v\n", coffees)))
		if err != nil {
			cfg.Logger.Println(err)
		}
	})

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
