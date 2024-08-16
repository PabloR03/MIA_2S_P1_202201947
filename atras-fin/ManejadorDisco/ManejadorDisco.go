package ManejadorDisco

import (
	"encoding/binary"
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
	// Validar el tamaño (size)
	if size <= 0 {
		fmt.Println("Error: El tamaño debe ser mayor que 0.")
		return
	}
	// Validar el ajuste (fit)
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: El ajuste debe ser BF, WF, o FF.")
		return
	}
	// Validar la unidad (unit)
	if unit != "k" && unit != "m" {
		fmt.Println("Error: La unidad debe ser K o M.")
		return
	}
	// Validar la ruta (path)
	if path == "" {
		fmt.Println("Error: La ruta es obligatoria.")
		return
	}

	// Crear el archivo en la ruta especificada
	err := Utilidades.CreateFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	// Convertir el tamaño a bytes
	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}
	// Abrir el archivo para escritura
	archivo, err := Utilidades.OpenFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	// Inicializar el archivo con ceros
	for i := 0; i < size; i++ {
		err := Utilidades.WriteObject(archivo, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	// Inicializar el MBR
	var nuevo_mbr Estructura.MRB
	nuevo_mbr.MbrSize = int32(size)
	nuevo_mbr.Signature = rand.Int31()
	currentTime := time.Now()
	fechaFormateada := currentTime.Format("2006-01-02")
	copy(nuevo_mbr.CreationDate[:], fechaFormateada)
	copy(nuevo_mbr.Fit[:], fit)
	// Escribir el MBR en el archivo
	if err := Utilidades.WriteObject(archivo, nuevo_mbr, 0); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer archivo.Close()
	fmt.Println("Disco creado con éxito en la ruta: ", path)
	fmt.Println("======End MKDISK======")

}

func Rmdisk(path string) {
	fmt.Println("======INICIO RMDISK======")
	err := Utilidades.EliminarArchivo(path)
	if err != nil {
		fmt.Println("Error RMDISK: ", err)
		return
	}
	if path == "" {
		fmt.Println("Error RMDISK: La ruta es obligatoria.")
		return
	}

	fmt.Println("Disco eliminado con éxito en la ruta: ", path)
	fmt.Println("======End RMDISK======")
}

func Fdisk(size int, unit string, path string, type_ string, fit string, name string) {
	fmt.Println("======Start FDISK======")
	fmt.Println("-------------------------------------------------------------")
	// Validar el tamaño (size)
	if size <= 0 {
		fmt.Println("Error: Tamaño debe ser mayor que 0.")
		return
	}

	// Validar la unidad (unit)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unidad debe ser B, K, M.")
		return
	}

	// Validar la ruta (path)
	if path == "" {
		fmt.Println("Error: La ruta es obligatoria.")
		return
	}

	// Validar el tipo (type)
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Tipo debe ser P, E, L.")
		return
	}

	// Validar el ajuste (fit)
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Ajuste debe ser BF, WF o FF")
		return
	}

	// Validar el nombre (name)
	if name == "" {
		fmt.Println("Error: El nombre es obligatorio.")
		return
	}

	// Convertir el tamaño a bytes
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// Abrir archivo binario
	archivo, err := Utilidades.OpenFile(path)
	if err != nil {
		return
	}

	var MBRTemporalDisco Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &MBRTemporalDisco, 0); err != nil {
		return
	}

	// Calcular el espacio restante
	espacioUsado := int32(0)
	for i := 0; i < 4; i++ {
		espacioUsado += MBRTemporalDisco.Partitions[i].Size
	}
	espacioRestante := MBRTemporalDisco.MbrSize - espacioUsado

	// Validar que el tamaño de la nueva partición no exceda el espacio restante
	if int32(size) > espacioRestante {
		fmt.Println("Error: El tamaño de la partición excede el espacio disponible.")
		fmt.Println("Tamaño restante disponible:", espacioRestante, "bytes")
		return
	}

	// Aquí continuarías con la lógica para agregar la partición, como antes
	var contador = 0
	var vacio = int32(0)

	for i := 0; i < 4; i++ {
		if MBRTemporalDisco.Partitions[i].Size != 0 {
			contador++
			vacio = MBRTemporalDisco.Partitions[i].Start + MBRTemporalDisco.Partitions[i].Size
		}
	}

	for i := 0; i < 4; i++ {
		if MBRTemporalDisco.Partitions[i].Size == 0 {
			MBRTemporalDisco.Partitions[i].Size = int32(size)
			if contador == 0 {
				MBRTemporalDisco.Partitions[i].Start = int32(binary.Size(MBRTemporalDisco))
			} else {
				MBRTemporalDisco.Partitions[i].Start = vacio
			}
			copy(MBRTemporalDisco.Partitions[i].Name[:], name)
			copy(MBRTemporalDisco.Partitions[i].Fit[:], fit)
			copy(MBRTemporalDisco.Partitions[i].Status[:], "0")
			copy(MBRTemporalDisco.Partitions[i].Type[:], type_)
			MBRTemporalDisco.Partitions[i].Correlative = int32(contador + 1)
			break
		}
	}

	if err := Utilidades.WriteObject(archivo, MBRTemporalDisco, 0); err != nil {
		return
	}

	var TempMBR2 Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &TempMBR2, 0); err != nil {
		return
	}
	Estructura.PrintMBR(TempMBR2)

	defer archivo.Close()

	fmt.Println("------------------")
	fmt.Println("Tamaño del disco:", MBRTemporalDisco.MbrSize, "bytes")
	fmt.Println("Tamaño utilizado:", espacioUsado, "bytes")
	fmt.Println("Tamaño restante:", espacioRestante, "bytes")
	fmt.Println("------------------")

	fmt.Println("Partición creada con éxito en la ruta:", path)
}
