package main

import (
	"github.com/nats-io/nats.go"
	"net/http"
	"valantonini/go-coffee-service/coffee-service/config"
	"valantonini/go-coffee-service/coffee-service/service"
)

func main() {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		panic("unable to create configuration")
	}

	cfg.Logger.Printf("connecting to nats on %v\n", cfg.NatsAddress)
	nc, err := nats.Connect(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Fatal(err.Error())
	}
	cfg.Logger.Println("connected to nats")

	cfg.Logger.Println("registering coffee services")
	coffeeService := service.NewCoffeeService(cfg, nc)
	http.HandleFunc("/coffees", coffeeService.List)
	http.HandleFunc("/coffee/add", coffeeService.Add)
	http.Handle("/health", service.NewHealth(cfg.Logger))
	cfg.Logger.Println("coffee services registered")

	cfg.Logger.Printf("starting server on %v", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Fatal(err)
	}
}
