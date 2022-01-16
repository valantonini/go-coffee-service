package main

import (
	"fmt"
	"log"
	"net/http"
	"valantonini/go-coffee-service/coffee-service/data"
)

func main() {
	repository, _ := data.InitRepository()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, _ := repository.Find()

		res, err := data.ToJSON()
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
