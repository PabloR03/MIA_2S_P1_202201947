package Analizador

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"proyecto1/ManejadorDisco"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

// Estructura para manejar la entrada JSON
type CommandRequest struct {
	Text string `json:"text"`
}

// Estructura para la respuesta JSON
type CommandResponse struct {
	Output string `json:"text"`
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

		command, params := getCommandAndParams(request.Text)

		output := AnalyzeCommand(command, params)

		response := CommandResponse{Output: output}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Obtiene el comando y los parámetros
func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

func AnalyzeCommand(command string, params string) string {
	var output strings.Builder
	if strings.Contains(command, "mkdisk") {
		output.WriteString(FuncionMkdisk(params))
	}
	if strings.Contains(command, "rmdisk") {
		output.WriteString(FuncionRmdisk(params))
	}
	return output.String()
}

func FuncionMkdisk(params string) string {
	var output strings.Builder
	// Definir flag
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamano")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	// Parse flag
	fs.Parse(strings.Fields(params))

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path":
			fs.Set(flagName, flagValue)
		default:
			output.WriteString("Error: Flag not found\n")
		}
	}

	// Validaciones
	if *size <= 0 {
		output.WriteString("Error: Size must be greater than 0\n")
		return output.String()
	}

	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		output.WriteString("Error: Fit must be 'bf', 'ff', or 'wf'\n")
		return output.String()
	}

	if *unit != "k" && *unit != "m" {
		output.WriteString("Error: Unit must be 'k' or 'm'\n")
		return output.String()
	}

	if *path == "" {
		output.WriteString("Error: Path is required\n")
		return output.String()
	}

	// Llamamos a la funcion
	ManejadorDisco.Mkdisk(*size, *fit, *unit, *path)
	return output.String()
}

func FuncionRmdisk(params string) string {
	var output strings.Builder
	// Preguntar si desea eliminar el disco
	output.WriteString("¿Desea eliminar el disco? (s/n): ")
	var respuesta string
	fmt.Scanln(&respuesta)
	respuesta = strings.ToLower(respuesta)

	if respuesta == "s" {
		// Definir flag
		fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
		path := fs.String("path", "", "Ruta")

		// Parse flag
		fs.Parse(strings.Fields(params))

		// Encontrar la flag en el input
		matches := re.FindAllStringSubmatch(params, -1)

		// Process the input
		for _, match := range matches {
			flagName := strings.ToLower(match[1])
			flagValue := strings.ToLower(match[2])
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "path":
				fs.Set(flagName, flagValue)
			default:
				output.WriteString("Error: Flag not found\n")
			}
		}

		// Validaciones
		if *path == "" {
			output.WriteString("Error: Path is required\n")
			return output.String()
		}

		// Llamamos a la funcion
		ManejadorDisco.Rmdisk(*path)
	} else if respuesta == "n" {
		output.WriteString("No se elimino el disco\n")
	}
	return output.String()
}

func main() {
	http.HandleFunc("/api/message", commandHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
