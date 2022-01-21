package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"valantonini/go-coffee-service/coffee-service/config"
	"valantonini/go-coffee-service/coffee-service/data"
	"valantonini/go-coffee-service/coffee-service/events"
	"valantonini/go-coffee-service/coffee-service/service"
)

func main() {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		panic("unable to create configuration")
	}
	repository, _ := data.InitRepository()

	cfg.Logger.Printf("connecting to nats on %v\n", cfg.NatsAddress)
	nc, err := nats.Connect(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Fatal(err.Error())
	}
	cfg.Logger.Println("connected to nats")

	coffeeService := service.NewCoffeeService(cfg)
	http.HandleFunc("/list", coffeeService.List)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			printRoutes(w, r)
			return
		}

		var requestData map[string]string
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&requestData)
		newItem := repository.Add(requestData["name"])

		res, err := json.Marshal(newItem)
		if err != nil {
			log.Printf("Error during JSON marshal. Err: %s", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		nc.Publish(events.CoffeeAdded, res)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	http.Handle("/health", service.NewHealth(cfg.Logger))

	cfg.Logger.Printf("starting server on %v", cfg.BindAddress)

	if err := http.ListenAndServe(cfg.BindAddress, nil); err != nil {
		log.Fatal(err)
	}
}
func printRoutes(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"GET /list": "retrieves a list of all coffees",
		"POST /add": "adds a coffee",
	}

	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error during JSON marshal. Err: %s", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
