package Analizador

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"proyecto1/ManejadorDisco"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func Analizar(texto string) string {
	var buffer bytes.Buffer

	scanner := bufio.NewScanner(strings.NewReader(texto))
	for scanner.Scan() {
		entrada := scanner.Text()
		if len(entrada) == 0 || entrada[0] == '#' {
			continue
		}
		entrada = strings.TrimSpace(entrada)
		command, params := getCommandAndParams(entrada)
		fmt.Fprintln(&buffer, "Comando:", command, "Parametros:", params)
		AnalyzeCommnad(command, params, &buffer)
	}

	return buffer.String()
}

func AnalyzeCommnad(command string, params string, buffer *bytes.Buffer) {
	// Pasa el buffer a las funciones
	if strings.Contains(command, "mkdisk") {
		Funcion_mkdisk(params, buffer)
	} else if strings.Contains(command, "fdisk") {
		Funcion_fdisk(params, buffer)
	} else if strings.Contains(command, "rmdisk") {
		Funcion_rmdisk(params, buffer)
	} else if strings.Contains(command, "mount") {
		Funcion_mount(params, buffer)
	}
}

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

func Funcion_mkdisk(params string, writer io.Writer) {
	// Define flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamano")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	// Parse los argumentos desde params en lugar de os.Args
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])
		valorFlag = strings.Trim(valorFlag, "\"")
		switch nombreFlag {
		case "size", "fit", "unit", "path":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Fprintf(writer, "Error: Parámetro no encontrado.\n")
		}
	}

	fs.Parse([]string{}) // Asegúrate de no intentar parsear argumentos adicionales de os.Args

	ManejadorDisco.Mkdisk(*size, *fit, *unit, *path, writer.(*bytes.Buffer))
}

func Funcion_rmdisk(params string, writer io.Writer) {
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
			fmt.Fprintf(writer, "Error: Parámetro no encontrado.\n")
		}
	}
	ManejadorDisco.Rmdisk(*path, writer.(*bytes.Buffer))
}

func Funcion_fdisk(input string, writer io.Writer) {
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
			fmt.Fprintf(writer, "Error: Parámetro no encontrado.\n")
		}
	}
	ManejadorDisco.Fdisk(*size, *unit, *path, *type_, *fit, *name, writer.(*bytes.Buffer))
}

func Funcion_mount(input string, writer io.Writer) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	ruta := fs.String("path", "", "Ruta")
	nombre := fs.String("name", "", "Nombre")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		nombreFlag := match[1]
		valorFlag := strings.ToLower(match[2])

		valorFlag = strings.Trim(valorFlag, "\"")

		switch nombreFlag {
		case "path", "name":
			fs.Set(nombreFlag, valorFlag)
		default:
			fmt.Println("Error: Parámetro no encontrado.")
		}
	}
	ManejadorDisco.Mount(*ruta, *nombre, writer.(*bytes.Buffer))
}
