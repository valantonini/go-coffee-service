package main

import (
	"github.com/nats-io/nats.go"
	"net/http"
	"valantonini/go-coffee-service/product-service/config"
	"valantonini/go-coffee-service/product-service/data"
	"valantonini/go-coffee-service/product-service/service"
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

	cfg.Logger.Println("registering product service handlers")
	productService := service.NewCoffeeService(repo, nc, cfg.Logger)
	http.HandleFunc("/coffees", productService.List)
	http.HandleFunc("/coffee/add", productService.Add)
	http.Handle("/health", service.NewHealth(cfg.Logger))
	cfg.Logger.Println("product services registered")

	cfg.Logger.Printf("starting server on %v", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Fatal(err)
	}
}
