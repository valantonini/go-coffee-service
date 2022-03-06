package main

import (
	"github.com/gorilla/mux"
	"github.com/valantonini/go-coffee-service/cmd/product-service/data"
	"github.com/valantonini/go-coffee-service/cmd/product-service/events"
	"github.com/valantonini/go-coffee-service/cmd/product-service/service"
	"github.com/valantonini/go-coffee-service/internal/pkg/config"
	"github.com/valantonini/go-coffee-service/internal/pkg/health"
	"net/http"
	"time"
)

func setContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.NewConfigFromEnv("product-service")

	cfg.Logger.Info("connecting to nats", "natAddress", cfg.NatsAddress)
	nc, err := events.NewNatsConnection(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Error(err.Error())
	}
	defer nc.Close()
	cfg.Logger.Info("connected to nats")

	cfg.Logger.Info("initialising db connection")
	db, err := data.NewDbConnection()
	if err != nil {
		cfg.Logger.Error(err.Error())
	}
	cfg.Logger.Info("db connection initialised")

	cfg.Logger.Info("initialising test data")
	err = data.InitTestData(db)
	if err != nil {
		cfg.Logger.Error(err.Error())
	}
	cfg.Logger.Info("test data initialised")

	cfg.Logger.Info("initialising repository")
	repo, err := data.NewMongoCoffeeRepository(db)
	if err != nil {
		cfg.Logger.Error("unable to initialise repo")
	}
	cfg.Logger.Info("repository initialised")

	cfg.Logger.Info("registering product service consumers")
	handlerService := service.NewConsumerService(repo, nc, cfg.Logger)
	handlerService.RegisterConsumer("get-coffees", handlerService.GetCoffees)
	cfg.Logger.Info("product service consumers registered")

	cfg.Logger.Info("registering outbox")
	outboxRepo, _ := data.NewMongoOutboxRepository(db)
	outbox := service.NewOutbox(&outboxRepo, nc)
	cancelOutbox := outbox.StartBackgroundPolling(500 * time.Millisecond)
	defer cancelOutbox()
	cfg.Logger.Info("outbox registered")

	cfg.Logger.Info("registering product service http handlers")
	r := mux.NewRouter()
	r.Use(setContentTypeMiddleware)
	productService := service.NewCoffeeService(repo, &outbox, cfg.Logger)
	productService.RegisterRoutes(r)
	r.Handle("/health", health.NewHealthService(cfg.Logger)).Methods(http.MethodGet)
	http.Handle("/", r)
	cfg.Logger.Info("product service http handlers registered")

	cfg.Logger.Info("starting server on", "bindAddress", cfg.BindAddress)
	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		cfg.Logger.Error(err.Error())
	}
}
