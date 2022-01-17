package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"valantonini/go-coffee-service/coffee-service/data"
)

func main() {
	repository, _ := data.InitRepository()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		printRoutes(w, r)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		data := repository.Find()

		res, err := data.ToJSON()
		if err != nil {
			log.Printf("Error during JSON marshal. Err: %s", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
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
