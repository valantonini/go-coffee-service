package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/valantonini/go-coffee-service/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfigFromEnv()

	cfg.Logger.Printf("connecting to nats on %v\n", cfg.NatsAddress)
	nc, err := nats.Connect(cfg.NatsAddress)
	if err != nil {
		cfg.Logger.Fatal(err.Error())
	}
	defer nc.Close()
	cfg.Logger.Println("connected to nats")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("%v %v", cfg.BindAddress, cfg.NatsAddress)))
	})
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
