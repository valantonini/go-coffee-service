package main

import (
	"github.com/nats-io/nats.go"
	"net/http"
	"valantonini/go-coffee-service/coffee-service/config"
	"valantonini/go-coffee-service/coffee-service/data"
	"valantonini/go-coffee-service/coffee-service/service"
)

func main() {
	cfg := config.NewConfigFromEnv()

	cfg.Logger.Printf("connecting to nats on %v\n", cfg.NatsAddress)
	nc, err := nats.Connect(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Fatal(err.Error())
	}
	cfg.Logger.Println("connected to nats")

	cfg.Logger.Print("initialising repository")
	repo, err := data.InitRepository()
	if err != nil {
		cfg.Logger.Fatalf("unable to initialise repo")
	}
	cfg.Logger.Print("repository initialised")

	cfg.Logger.Println("registering coffee service handlers")
	coffeeService := service.NewCoffeeService(repo, nc, cfg.Logger)
	http.HandleFunc("/coffees", coffeeService.List)
	http.HandleFunc("/coffee/add", coffeeService.Add)
	http.Handle("/health", service.NewHealth(cfg.Logger))
	cfg.Logger.Println("coffee services registered")

	cfg.Logger.Printf("starting server on %v", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Fatal(err)
	}
}
