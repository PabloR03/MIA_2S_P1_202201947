package Utilidades

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// Funcion para crear un archivo binario
func CreateFile(name string) error {
	// Asignar directorio
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Crear archivo
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funcion para abrir el archivobinario en Lectura/Escritura
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

func EliminarArchivo(nombre string) error {
	if _, err := os.Stat(nombre); os.IsNotExist(err) {
		fmt.Println("Error: El archivo no existe.")
		return err
	}
	err := os.Remove(nombre)
	if err != nil {
		fmt.Println("Error al eliminar el archivo: ", err)
		return err
	}
	return nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}
