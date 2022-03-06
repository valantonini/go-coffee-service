package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/cmd/order-service/gateway"
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
	cfg := config.NewConfigFromEnv("order-service")

	cfg.Logger.Info("connecting to nats", "natAddress", cfg.NatsAddress)
	// var b gateway.Bus
	nc, err := nats.Connect(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Error(err.Error())
	}
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		cfg.Logger.Error(err.Error())
	}
	cfg.Logger.Info("connected to nats")

	cfg.Logger.Info("retrieving coffee from coffee-service")
	coffeeService := gateway.NewCoffeeServiceGateway(ec)
	coffees, err := coffeeService.GetAll()
	if err != nil {
		cfg.Logger.Error(err.Error())
	}
	cfg.Logger.Info("coffees retrieved", "count", len(coffees))

	cfg.Logger.Info("registering order service http handlers")
	r := mux.NewRouter()
	r.Use(setContentTypeMiddleware)
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		err := json.NewEncoder(writer).Encode(coffees)
		if err != nil {
			cfg.Logger.Error(err.Error())
		}
	})
	r.Handle("/health", health.NewHealthService(cfg.Logger)).Methods(http.MethodGet)
	http.Handle("/", r)
	cfg.Logger.Info("order service http handlers registered")

	cfg.Logger.Info("starting server", "binAddress", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Error(err.Error())
	}
}
