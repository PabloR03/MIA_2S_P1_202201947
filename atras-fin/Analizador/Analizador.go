package Analizador

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"proyecto1/ManejadorDisco"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

// Analiza el texto de entrada
func Analizar(texto string) {
	scanner := bufio.NewScanner(strings.NewReader(texto))
	for scanner.Scan() {
		entrada := scanner.Text()
		if len(entrada) == 0 || entrada[0] == '#' {
			continue
		}
		entrada = strings.TrimSpace(entrada)
		command, params := getCommandAndParams(entrada)
		fmt.Println("Comando:", command, "Parametros:", params)
		AnalyzeCommnad(command, params)
	}
}

// Obtiene el comando y los parametros
func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		for i := 1; i < len(parts); i++ {
			parts[i] = strings.ToLower(parts[i])
		}
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

func AnalyzeCommnad(command string, params string) {

	if strings.Contains(command, "mkdisk") {
		Funcion_mkdisk(params)
	} else if strings.Contains(command, "fdisk") {
		Funcion_fdisk(params)
	} else if strings.Contains(command, "rmdisk") {
		Funcion_rmdisk(params)
	}
	/* else if strings.Contains(command, "mount") {
		fn_mount(params)
	} else if strings.Contains(command, "mkfs") {
		fn_mkfs(params)
	} else if strings.Contains(command, "login") {
		fn_login(params)
	} else if strings.Contains(command, "logout") {
		fn_logout()
	} else if strings.Contains(command, "mkusr") {
		fn_mkusr(params)
	} else {
		fmt.Println("Error: Command not found")
	} */

}

func Funcion_mkdisk(params string) {
	// Define flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamano")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])
		valorFlag = strings.Trim(valorFlag, "\"")
		switch nombreFlag {
		case "size", "fit", "unit", "path":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Println("Error: Parámetro no encontrado.")
		}
	}
	ManejadorDisco.Mkdisk(*size, *fit, *unit, *path)
}

func Funcion_rmdisk(params string) {
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])
		valorFlag = strings.Trim(valorFlag, "\"")
		switch nombreFlag {
		case "path":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Println("Error: Parámetro no encontrado.")
		}
	}
	ManejadorDisco.Rmdisk(*path)
}

func Funcion_fdisk(input string) {
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	unit := fs.String("unit", "k", "Unidad")
	path := fs.String("path", "", "Ruta")
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "wf", "Ajuste")
	name := fs.String("name", "", "Nombre")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])

		valorFlag = strings.Trim(valorFlag, "\"")

		switch nombreFlag {
		case "size", "unit", "path", "type", "fit", "name":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Println("Error: Parámetro no encontrado.")
		}
	}
	ManejadorDisco.Fdisk(*size, *unit, *path, *type_, *fit, *name)
}
