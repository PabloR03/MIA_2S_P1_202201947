package ManejadorDisco

import (
	"fmt"
	"math/rand"
	"proyecto1/Estructura"
	"proyecto1/Utilidades"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string) {
	fmt.Println("======INICIO MKDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Fit:", fit)
	fmt.Println("Unit:", unit)
	fmt.Println("Path:", path)

	// Validar fit bf/ff/wf
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Fit debe ser bf, wf or ff")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: Size debe ser mayo a  0")
		return
	}

	// Validar unidar k - m
	if unit != "k" && unit != "m" {
		fmt.Println("Error: Las unidades validas son k o m")
		return
	}

	// Create file
	err := Utilidades.CreateFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	/*
		Si el usuario especifica unit = "k" (Kilobytes), el tamaño se multiplica por 1024 para convertirlo a bytes.
		Si el usuario especifica unit = "m" (Megabytes), el tamaño se multiplica por 1024 * 1024 para convertirlo a MEGA bytes.
	*/
	// Asignar tamanio
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// Open bin file
	file, err := Utilidades.OpenFile(path)
	if err != nil {
		return
	}

	// Escribir los 0 en el archivo

	// create array of byte(0)
	for i := 0; i < size; i++ {
		err := Utilidades.WriteObject(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// Create a new instance of MRB
	var newMRB Estructura.MRB
	newMRB.MbrSize = int32(size)
	newMRB.Signature = rand.Int31()
	copy(newMRB.Fit[:], fit)
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	copy(newMRB.CreationDate[:], formattedDate)

	// Write object in bin file
	if err := Utilidades.WriteObject(file, newMRB, 0); err != nil {
		return
	}

	var TempMBR Estructura.MRB
	// Read object from bin file
	if err := Utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	Estructura.PrintMBR(TempMBR)

	// Close bin file
	defer file.Close()

	fmt.Println("======End MKDISK======")

}

func Rmdisk(path string) {
	fmt.Println("======INICIO RMDISK======")
	fmt.Println("Path:", path)

	// Validar path
	if path == "" {
		fmt.Println("Error: Path es requerido")
		return
	}

	// Eliminar archivo
	err := Utilidades.DeleteFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("======End RMDISK======")
}
