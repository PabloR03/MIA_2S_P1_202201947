package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	http.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			log.Println("Received GET request")
			message := Message{Text: "Hello from Go!"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)
			return
		}

		if r.Method == http.MethodPost {
			log.Println("Received POST request")
			var message Message
			if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
				log.Println("Error decoding JSON:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			log.Println("Received POST request with message:", message.Text)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	})

	// Configura CORS para permitir solicitudes desde tu frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
