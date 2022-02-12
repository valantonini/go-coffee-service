package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/internal/pkg/config"
	"github.com/valantonini/go-coffee-service/order-service/gateway"
	"github.com/valantonini/go-coffee-service/order-service/service"
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
	cfg.Logger.Printf("%v coffees retrieved", len(coffees))

	cfg.Logger.Println("registering order service http handlers")
	r := mux.NewRouter()
	r.Use(setContentTypeMiddleware)
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		err := json.NewEncoder(writer).Encode(coffees)
		if err != nil {
			cfg.Logger.Println(err)
		}
	})
	r.Handle("/health", service.NewHealth(cfg.Logger))
	http.Handle("/", r)
	cfg.Logger.Println("order service http handlers registered")

	cfg.Logger.Printf("starting server on %v", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Fatal(err)
	}
}
