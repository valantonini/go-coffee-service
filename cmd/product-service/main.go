package main

import (
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"github.com/valantonini/go-coffee-service/cmd/product-service/service"
	"github.com/valantonini/go-coffee-service/internal/pkg/config"
	"github.com/valantonini/go-coffee-service/internal/pkg/health"
	"net/http"
)

func setContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.NewConfigFromEnv()

	cfg.Logger.Printf("connecting to nats on %v\n", cfg.NatsAddress)
	nc, err := nats.Connect(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Fatal(err.Error())
	}
	defer nc.Close()
	cfg.Logger.Println("connected to nats")

	cfg.Logger.Print("initialising repository")
	repo, err := data.InitInMemoryRepository()
	if err != nil {
		cfg.Logger.Fatalf("unable to initialise repo")
	}
	cfg.Logger.Print("repository initialised")

	cfg.Logger.Println("registering product service consumers")
	handlerService := service.NewConsumerService(repo, nc, cfg.Logger)
	handlerService.RegisterConsumer("get-coffees", handlerService.GetCoffees)
	cfg.Logger.Println("product service consumers registered")

	cfg.Logger.Println("registering product service http handlers")
	r := mux.NewRouter()
	r.Use(setContentTypeMiddleware)
	productService := service.NewCoffeeService(repo, nc, cfg.Logger)
	productService.RegisterRoutes(r)
	r.Handle("/health", health.NewHealthService(cfg.Logger)).Methods(http.MethodGet)
	http.Handle("/", r)
	cfg.Logger.Println("product service http handlers registered")

	cfg.Logger.Printf("starting server on %v", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Fatal(err)
	}
}
