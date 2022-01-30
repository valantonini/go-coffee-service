package main

import (
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/config"
	"github.com/valantonini/go-coffee-service/product-service/data"
	"github.com/valantonini/go-coffee-service/product-service/service"
	"net/http"
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
	repo, err := data.InitInMemoryRepository()
	if err != nil {
		cfg.Logger.Fatalf("unable to initialise repo")
	}
	cfg.Logger.Print("repository initialised")

	cfg.Logger.Println("registering product service handlers")
	productService := service.NewCoffeeService(repo, nc, cfg.Logger)
	r := mux.NewRouter()
	r.HandleFunc("/coffees", productService.List)
	r.HandleFunc("/coffee/add", productService.Add)
	r.HandleFunc("/coffee/{id}", productService.Get)
	r.Handle("/health", service.NewHealth(cfg.Logger))
	http.Handle("/", r)
	cfg.Logger.Println("product services registered")

	cfg.Logger.Printf("starting server on %v", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Fatal(err)
	}
}
