package ManejadorDisco

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"proyecto1/Estructura"
	"proyecto1/Utilidades"
	"time"
)

func Mkdisk(size int, fit string, unit string, path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "======INICIO MKDISK======")
	fmt.Fprintln(buffer, "Size:", size)
	fmt.Fprintln(buffer, "Fit:", fit)
	fmt.Fprintln(buffer, "Unit:", unit)
	fmt.Fprintln(buffer, "Path:", path)

	// Validar el tamaño (size)
	if size <= 0 {
		fmt.Fprintln(buffer, "Error: El tamaño debe ser mayor que 0.")
		return
	}

	// Validar el ajuste (fit)
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Fprintln(buffer, "Error: El ajuste debe ser BF, WF, o FF.")
		return
	}

	// Validar la unidad (unit)
	if unit != "k" && unit != "m" {
		fmt.Fprintln(buffer, "Error: La unidad debe ser K o M.")
		return
	}

	// Validar la ruta (path)
	if path == "" {
		fmt.Fprintln(buffer, "Error: La ruta es obligatoria.")
		return
	}

	// Crear el archivo en la ruta especificada
	err := Utilidades.CreateFile(path)
	if err != nil {
		fmt.Fprintln(buffer, "Error: ", err)
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
		fmt.Fprintln(buffer, "Error: ", err)
		return
	}

	// Inicializar el archivo con ceros
	for i := 0; i < size; i++ {
		err := Utilidades.WriteObject(archivo, byte(0), int64(i))
		if err != nil {
			fmt.Fprintln(buffer, "Error: ", err)
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
		fmt.Fprintln(buffer, "Error: ", err)
		return
	}
	defer archivo.Close()

	fmt.Fprintln(buffer, "Disco creado con éxito en la ruta: ", path)
	fmt.Fprintln(buffer, "======End MKDISK======")
}

func Rmdisk(path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "======INICIO RMDISK======")

	// Validar la ruta (path)
	if path == "" {
		fmt.Fprintln(buffer, "Error RMDISK: La ruta es obligatoria.")
		return
	}

	// Eliminar el archivo en la ruta especificada
	err := Utilidades.EliminarArchivo(path)
	if err != nil {
		fmt.Fprintln(buffer, "Error RMDISK:", err)
		return
	}

	fmt.Fprintln(buffer, "Disco eliminado con éxito en la ruta:", path)
	fmt.Fprintln(buffer, "======End RMDISK======")
}

func Fdisk(size int, unit string, path string, type_ string, fit string, name string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "======Start FDISK======")
	fmt.Fprintln(buffer, "-------------------------------------------------------------")
	fmt.Fprintln(buffer, "Size:", size)
	fmt.Fprintln(buffer, "Unit:", unit)
	fmt.Fprintln(buffer, "Path:", path)
	fmt.Fprintln(buffer, "Type:", type_)
	fmt.Fprintln(buffer, "Fit:", fit)
	fmt.Fprintln(buffer, "Name:", name)

	// Validar el tamaño (size)
	if size <= 0 {
		fmt.Fprintln(buffer, "Error: Tamaño debe ser mayor que 0.")
		return
	}

	// Validar la unidad (unit)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Fprintln(buffer, "Error: Unidad debe ser B, K, M.")
		return
	}

	// Validar la ruta (path)
	if path == "" {
		fmt.Fprintln(buffer, "Error: La ruta es obligatoria.")
		return
	}

	// Validar el tipo (type)
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Fprintln(buffer, "Error: Tipo debe ser P, E, L.")
		return
	}

	// Validar el ajuste (fit)
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Fprintln(buffer, "Error: Ajuste debe ser BF, WF o FF")
		return
	}

	// Validar el nombre (name)
	if name == "" {
		fmt.Fprintln(buffer, "Error: El nombre es obligatorio.")
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
		fmt.Fprintln(buffer, "Error al abrir el archivo:", err)
		return
	}

	var MBRTemporalDisco Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &MBRTemporalDisco, 0); err != nil {
		fmt.Fprintln(buffer, "Error al leer el MBR:", err)
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
		fmt.Fprintln(buffer, "Error: El tamaño de la partición excede el espacio disponible.")
		fmt.Fprintln(buffer, "Tamaño restante disponible:", espacioRestante, "bytes")
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
		fmt.Fprintln(buffer, "Error al escribir el MBR:", err)
		return
	}

	var TempMBR2 Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &TempMBR2, 0); err != nil {
		fmt.Fprintln(buffer, "Error al leer el MBR actualizado:", err)
		return
	}
	Estructura.PrintMBR(TempMBR2)

	defer archivo.Close()

	fmt.Fprintln(buffer, "------------------")
	fmt.Fprintln(buffer, "Tamaño del disco:", MBRTemporalDisco.MbrSize, "bytes")
	fmt.Fprintln(buffer, "Tamaño utilizado:", espacioUsado, "bytes")
	fmt.Fprintln(buffer, "Tamaño restante:", espacioRestante, "bytes")
	fmt.Fprintln(buffer, "------------------")

	fmt.Fprintln(buffer, "Partición creada con éxito en la ruta:", path)
	fmt.Fprintln(buffer, "======End FDISK======")
}
