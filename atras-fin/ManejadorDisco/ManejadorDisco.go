package ManejadorDisco

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"proyecto1/Estructura"
	"proyecto1/Utilidades"
	"strings"
	"time"
)

var RegistroDisco []Estructura.MountId
var LetraInicial = 64

func AgregarMountID(mountId Estructura.MountId) {
	RegistroDisco = append(RegistroDisco, mountId)
}

func NumeroLetraMountID(path string) (numero int32, letra int32) {
	for _, mount := range RegistroDisco {
		if mount.MIRuta == path {
			numero = mount.MINumero
			letra = mount.MILetra
			return numero, letra
		}
	}
	return 0, 0
}

func IncrementoNumeroMountID(path string) {
	for i, mount := range RegistroDisco {
		if mount.MIRuta == path {
			RegistroDisco[i].MINumero++
			break
		}
	}
}

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

	// Validar el fit (fit)
	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Fprintln(buffer, "Error: El fit debe ser BF, WF, o FF.")
		return
	}

	// Validar la unit (unit)
	if unit != "k" && unit != "m" {
		fmt.Fprintln(buffer, "Error: La unit debe ser K o M.")
		return
	}

	// Validar la path (path)
	if path == "" {
		fmt.Fprintln(buffer, "Error: La path es obligatoria.")
		return
	}

	// Crear el archivo en la path especificada
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

	LetraInicial++
	AgregarMountID(Estructura.MountId{
		MIRuta:   path,
		MINumero: 1,
		MILetra:  int32(LetraInicial),
	})

	fmt.Fprintln(buffer, "Disco creado con éxito en la path: ", path)
	fmt.Fprintln(buffer, "======End MKDISK======")
}

func Rmdisk(path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "======INICIO RMDISK======")

	// Validar la path (path)
	if path == "" {
		fmt.Fprintln(buffer, "Error RMDISK: La path es obligatoria.")
		return
	}

	// Eliminar el archivo en la path especificada
	err := Utilidades.EliminarArchivo(path)
	if err != nil {
		fmt.Fprintln(buffer, "Error RMDISK:", err)
		return
	}

	fmt.Fprintln(buffer, "Disco eliminado con éxito en la path:", path)
	fmt.Fprintln(buffer, "======End RMDISK======")
}

func Fdisk(size int, unit string, path string, type_ string, fit string, name string, buffer *bytes.Buffer) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Path:", path)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)

	// Validar el tamaño (size)
	if size <= 0 {
		fmt.Println("Error: Tamaño debe ser mayor que 0.")
		return
	}
	// Validar la unit (unit)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unidad debe ser B, K, M.")
		return
	}
	// Validar la path (path)
	if path == "" {
		fmt.Println("Error: La path es obligatoria.")
		return
	}
	// Validar el type_ (type)
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Tipo de partición debe ser P, E, L.")
		return
	}
	// Validar el fit (fit)
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

	var MBRTemporal Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &MBRTemporal, 0); err != nil {
		return
	}

	for i := 0; i < 4; i++ {
		if strings.Contains(string(MBRTemporal.Partitions[i].Name[:]), name) {
			fmt.Println("Error: El nombre: ", name, "ya está en uso en las particiones.")
			return
		}
	}

	var ContadorPrimaria, ContadorExtendida, TotalParticiones int
	var EspacioUtilizado int32 = 0

	for i := 0; i < 4; i++ {
		if MBRTemporal.Partitions[i].Size != 0 {
			TotalParticiones++
			EspacioUtilizado += MBRTemporal.Partitions[i].Size

			if MBRTemporal.Partitions[i].Type[0] == 'p' {
				ContadorPrimaria++
			} else if MBRTemporal.Partitions[i].Type[0] == 'e' {
				ContadorExtendida++
			}
		}
	}

	if TotalParticiones >= 4 && type_ != "l" {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		return
	}
	if type_ == "e" && ContadorExtendida > 0 {
		fmt.Println("Error: Solo se permite una partición extendida por disco.")
		return
	}
	if type_ == "l" && ContadorExtendida == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		return
	}
	if EspacioUtilizado+int32(size) > MBRTemporal.MbrSize {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return
	}

	var vacio int32 = int32(binary.Size(MBRTemporal))
	if TotalParticiones > 0 {
		vacio = MBRTemporal.Partitions[TotalParticiones-1].Start + MBRTemporal.Partitions[TotalParticiones-1].Size
	}

	for i := 0; i < 4; i++ {
		if MBRTemporal.Partitions[i].Size == 0 {
			if type_ == "p" || type_ == "e" {
				MBRTemporal.Partitions[i].Size = int32(size)
				MBRTemporal.Partitions[i].Size = vacio
				copy(MBRTemporal.Partitions[i].Name[:], []byte(name))
				copy(MBRTemporal.Partitions[i].Fit[:], fit)
				copy(MBRTemporal.Partitions[i].Status[:], "0")
				copy(MBRTemporal.Partitions[i].Type[:], type_)
				MBRTemporal.Partitions[i].Correlative = int32(TotalParticiones + 1)
				if type_ == "e" {
					EBRInicio := vacio
					EBRNuevo := Estructura.EBR{
						Part_Fit:   [1]byte{fit[0]},
						Part_Start: EBRInicio,
						Part_Size:  0,
						Part_Next:  -1,
					}
					copy(EBRNuevo.Part_Name[:], "")
					Utilidades.WriteObject(archivo, EBRNuevo, int64(EBRInicio))
				}
				fmt.Println(buffer, "Partición creada exitosamente en la path: ", path, " con el nombre: ", name)
				break
			}
		}
	}

	if type_ == "l" {
		var ParticionExtendida *Estructura.Partition
		for i := 0; i < 4; i++ {
			if MBRTemporal.Partitions[i].Type[0] == 'e' {
				ParticionExtendida = &MBRTemporal.Partitions[i]
				break
			}
		}
		if ParticionExtendida == nil {
			fmt.Println("Error: No se encontró una partición extendida para crear la partición lógica.")
			return
		}

		EBRPosterior := ParticionExtendida.Size
		var EBRUltimo Estructura.EBR
		for {
			Utilidades.ReadObject(archivo, &EBRUltimo, int64(EBRPosterior))
			if strings.Contains(string(EBRUltimo.Part_Name[:]), name) {
				fmt.Println("Error: El nombre: ", name, "ya está en uso en las particiones.")
				return
			}
			if EBRUltimo.Part_Next == -1 {
				break
			}
			EBRPosterior = EBRUltimo.Part_Next
		}

		var EBRNuevoPosterior int32
		if EBRUltimo.Part_Size == 0 {
			EBRNuevoPosterior = EBRPosterior
		} else {
			EBRNuevoPosterior = EBRUltimo.Part_Size + EBRUltimo.Part_Size
		}

		if EBRNuevoPosterior+int32(size)+int32(binary.Size(Estructura.EBR{})) > ParticionExtendida.Size+ParticionExtendida.Size {
			fmt.Println("Error: No hay suficiente espacio en la partición extendida para esta partición lógica.")
			return
		}

		if EBRUltimo.Part_Size != 0 {
			EBRUltimo.Part_Next = EBRNuevoPosterior
			Utilidades.WriteObject(archivo, EBRUltimo, int64(EBRPosterior))
		}

		newEBR := Estructura.EBR{
			Part_Fit:   [1]byte{fit[0]},
			Part_Start: EBRNuevoPosterior + int32(binary.Size(Estructura.EBR{})),
			Part_Size:  int32(size),
			Part_Next:  -1,
		}
		copy(newEBR.Part_Name[:], newEBR.Part_Name[:])
		Utilidades.WriteObject(archivo, newEBR, int64(EBRNuevoPosterior))
		fmt.Println("Partición lógica creada exitosamente en la path: ", path, " con el nombre: ", name)
	}

	if err := Utilidades.WriteObject(archivo, MBRTemporal, 0); err != nil {
		return
	}
	defer archivo.Close()
	fmt.Println("======FIN FDISK======")
}

func Mount(ruta string, nombre string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "======INICIO Mount======")
	// Validar la ruta (path)
	if ruta == "" {
		fmt.Println("Error: La ruta es obligatoria.")
		return
	}
	// Validar el nombre (name)
	if nombre == "" {
		fmt.Println("Error: El nombre es obligatorio.")
		return
	}
	// Abrir archivo binario
	archivo, err := Utilidades.OpenFile(ruta)
	if err != nil {
		return
	}
	var MBRTemporal Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &MBRTemporal, 0); err != nil {
		return
	}

	var ParticionExiste = false
	for i := 0; i < 4; i++ {
		if strings.Contains(string(MBRTemporal.Partitions[i].Name[:]), nombre) {
			ParticionExiste = true
			break
		}
	}

	if !ParticionExiste {
		fmt.Println("Error: No se encontró la partición con el nombre especificado.")
		return
	} else {
		numero, letra := NumeroLetraMountID(ruta)
		for i := 0; i < 4; i++ {
			if strings.Contains(string(MBRTemporal.Partitions[i].Type[:]), "e") || strings.Contains(string(MBRTemporal.Partitions[i].Type[:]), "l") {
				fmt.Println("Error: No se puede montar una partición extendida o lógica.")
				return
			} else if strings.Contains(string(MBRTemporal.Partitions[i].Name[:]), nombre) {
				MBRTemporal.Partitions[i].Status[0] = '1'
				MBRTemporal.Partitions[i].PartId = [4]byte{'4', '7', byte(rune(letra)), byte('0' + numero)}
				fmt.Println("Partición montada exitosamente en la ruta: ", ruta)
				break
			}
		}
		IncrementoNumeroMountID(ruta)
	}
	if err := Utilidades.WriteObject(archivo, MBRTemporal, 0); err != nil {
		return
	}

	var TempMBR2 Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &TempMBR2, 0); err != nil {
		return
	}
	fmt.Println("----MBR-------")
	Estructura.PrintMBR(TempMBR2, buffer)
	fmt.Println("----MBR-------")
	defer archivo.Close()
	fmt.Fprintln(buffer, "======end Mount======")
}

/*
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

	// Validar la unit (unit)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Fprintln(buffer, "Error: Unidad debe ser B, K, M.")
		return
	}

	// Validar la path (path)
	if path == "" {
		fmt.Fprintln(buffer, "Error: La path es obligatoria.")
		return
	}

	// Validar el type_ (type)
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Fprintln(buffer, "Error: Tipo debe ser P, E, L.")
		return
	}

	// Validar el fit (fit)
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
	Estructura.PrintMBR(TempMBR2, buffer)

	defer archivo.Close()

	fmt.Fprintln(buffer, "------------------")
	fmt.Fprintln(buffer, "Tamaño del disco:", MBRTemporalDisco.MbrSize, "bytes")
	fmt.Fprintln(buffer, "Tamaño utilizado:", espacioUsado, "bytes")
	fmt.Fprintln(buffer, "Tamaño restante:", espacioRestante, "bytes")
	fmt.Fprintln(buffer, "------------------")

	fmt.Fprintln(buffer, "Partición creada con éxito en la path:", path)
	fmt.Fprintln(buffer, "======End FDISK======")
}
*/
