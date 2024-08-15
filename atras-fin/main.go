package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"proyecto1/Analizador"
	"strings"
)

// Estructura para manejar la entrada JSON
type CommandRequest struct {
	Text string `json:"text"`
}

// Estructura para la respuesta JSON
type CommandResponse struct {
	Text string `json:"text"`
}

// Handler para recibir y procesar los comandos
func commandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var request CommandRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		command, params := Analizador.GetCommandAndParams(request.Text)
		var output strings.Builder
		oldStdout := os.Stdout
		os.Stdout = &output
		Analizador.AnalyzeCommand(command, params)
		os.Stdout = oldStdout

		response := CommandResponse{Text: output.String()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/api/message", commandHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
