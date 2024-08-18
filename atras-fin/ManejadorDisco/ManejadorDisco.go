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
var LetraInicial = 65

func AgregarMountID(mountId Estructura.MountId) {
	RegistroDisco = append(RegistroDisco, mountId)
}

func NumeroLetraMountID(path string) (numero int32, letra int32) {
	for _, mount := range RegistroDisco {
		if mount.MIDpath == path {
			numero = mount.MIDnumber
			letra = mount.MIDletter
			return numero, letra
		}
	}
	return 0, 0
}

func IncrementoNumeroMountID(path string) {
	for i, mount := range RegistroDisco {
		if mount.MIDpath == path {
			RegistroDisco[i].MIDnumber++
			break
		}
	}
}

func Mkdisk(size int, fit string, unit string, path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=======================================INICIO MKDISK=======================================")
	fmt.Fprintln(buffer, "Size:", size)
	fmt.Fprintln(buffer, "Fit:", fit)
	fmt.Fprintln(buffer, "Unit:", unit)
	fmt.Fprintln(buffer, "Path:", path)

	// ================================= VALIDACIONES =================================
	if size <= 0 {
		fmt.Fprintln(buffer, "Error: El tamaño debe ser mayor que 0.")
		return
	}

	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Fprintln(buffer, "Error: El fit debe ser BF, WF, o FF.")
		return
	}

	if unit != "k" && unit != "m" {
		fmt.Fprintln(buffer, "Error: La unit debe ser K o M.")
		return
	}

	if path == "" {
		fmt.Fprintln(buffer, "Error: La path es obligatoria.")
		return
	}

	err := Utilidades.CreateFile(path)

	if err != nil {
		fmt.Fprintln(buffer, "Error: ", err)
		return
	}

	if unit == "k" {
		size = size * 1024
	} else {
		size = size * 1024 * 1024
	}

	// ================================= ABRIR ARCHIVO =================================
	archivo, err := Utilidades.OpenFile(path)

	if err != nil {
		fmt.Fprintln(buffer, "Error: ", err)
		return
	}

	// ================================= inicializar el archivo con 0
	for i := 0; i < size; i++ {
		err := Utilidades.WriteObject(archivo, byte(0), int64(i))
		if err != nil {
			fmt.Fprintln(buffer, "Error: ", err)
			return
		}
	}

	// ================================= Inicializar el MBR
	var nuevo_mbr Estructura.MRB
	nuevo_mbr.MRBSize = int32(size)
	nuevo_mbr.MRBSignature = rand.Int31()
	currentTime := time.Now()
	fechaFormateada := currentTime.Format("2006-01-02")
	copy(nuevo_mbr.MRBCreationDate[:], fechaFormateada)
	copy(nuevo_mbr.MRBFit[:], fit)

	// ================================= Escribir el MBR en el archivo
	if err := Utilidades.WriteObject(archivo, nuevo_mbr, 0); err != nil {
		fmt.Fprintln(buffer, "Error: ", err)
		return
	}
	defer archivo.Close()

	AgregarMountID(Estructura.MountId{
		MIDpath:   path,
		MIDnumber: 1,
		MIDletter: int32(LetraInicial),
	})
	LetraInicial++

	fmt.Fprintln(buffer, "Disco creado con éxito en la path: ", path)
	fmt.Fprintln(buffer, "=======================================End MKDISK=======================================")
}

func Rmdisk(path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=======================================INICIO RMDISK=======================================")

	// ================================= Validar la path (path)
	if path == "" {
		fmt.Fprintln(buffer, "Error RMDISK: La path es obligatoria.")
		return
	}

	// ================================= Eliminar el archivo en la path especificada
	err := Utilidades.DeleteFile(path)
	if err != nil {
		fmt.Fprintln(buffer, "Error RMDISK:", err)
		return
	}

	fmt.Fprintln(buffer, "Disco eliminado con éxito en la path:", path)
	fmt.Fprintln(buffer, "=======================================End RMDISK=======================================")
}

func Fdisk(size int, unit string, path string, type_ string, fit string, name string, buffer *bytes.Buffer) {
	fmt.Println("======================================= Start FDISK =======================================")
	fmt.Fprintln(buffer, "Size:", size)
	fmt.Fprintln(buffer, "Path:", path)
	fmt.Fprintln(buffer, "Name:", name)
	fmt.Fprintln(buffer, "Unit:", unit)
	fmt.Fprintln(buffer, "Type:", type_)
	fmt.Fprintln(buffer, "Fit:", fit)

	// ================================= VALIDACIONES =================================
	if size <= 0 {
		fmt.Println("Error: Tamaño debe ser mayor que 0.")
		return
	}

	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unidad debe ser B, K, M.")
		return
	}

	if path == "" {
		fmt.Println("Error: La path es obligatoria.")
		return
	}

	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Println("Error: Tipo de partición debe ser P, E, L.")
		return
	}

	if fit != "bf" && fit != "wf" && fit != "ff" {
		fmt.Println("Error: Ajuste debe ser BF, WF o FF")
		return
	}

	if name == "" {
		fmt.Println("Error: El nombre es obligatorio.")
		return
	}

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// ================================= ABRIR ARCHIVO =================================
	archivo, err := Utilidades.OpenFile(path)
	if err != nil {
		return
	}

	// ================================= LEER MBR =================================
	var MBRTemporal Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &MBRTemporal, 0); err != nil {
		return
	}

	// ================================= VALIDAR NOMBRE =================================
	for i := 0; i < 4; i++ {
		if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Name[:]), name) {
			fmt.Println("Error: El nombre: ", name, "ya está en uso en las particiones.")
			return
		}
	}

	var ContadorPrimaria, ContadorExtendida, TotalParticiones int
	var EspacioUtilizado int32 = 0

	for i := 0; i < 4; i++ {
		if MBRTemporal.MRBPartitions[i].PART_Size != 0 {
			TotalParticiones++
			EspacioUtilizado += MBRTemporal.MRBPartitions[i].PART_Size
			if MBRTemporal.MRBPartitions[i].PART_Type[0] == 'p' {
				ContadorPrimaria++
			} else if MBRTemporal.MRBPartitions[i].PART_Type[0] == 'e' {
				ContadorExtendida++
			}
		}
	}

	// ================================= VALIDAR PARTICIONES =================================
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
	if EspacioUtilizado+int32(size) > MBRTemporal.MRBSize {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return
	}

	var vacio int32 = int32(binary.Size(MBRTemporal))
	if TotalParticiones > 0 {
		vacio = MBRTemporal.MRBPartitions[TotalParticiones-1].PART_Start + MBRTemporal.MRBPartitions[TotalParticiones-1].PART_Size
	}

	for i := 0; i < 4; i++ {
		if MBRTemporal.MRBPartitions[i].PART_Size == 0 {
			if type_ == "p" || type_ == "e" {
				MBRTemporal.MRBPartitions[i].PART_Size = int32(size)
				MBRTemporal.MRBPartitions[i].PART_Start = vacio
				copy(MBRTemporal.MRBPartitions[i].PART_Name[:], []byte(name))
				copy(MBRTemporal.MRBPartitions[i].PART_Fit[:], fit)
				copy(MBRTemporal.MRBPartitions[i].PART_Status[:], "0")
				copy(MBRTemporal.MRBPartitions[i].PART_Type[:], type_)
				MBRTemporal.MRBPartitions[i].PART_Correlative = int32(TotalParticiones + 1)
				if type_ == "e" {
					EBRInicio := vacio
					EBRNuevo := Estructura.EBR{
						ERBFit:   [1]byte{fit[0]},
						ERBStart: EBRInicio,
						ERBSize:  0,
						ERBNext:  -1,
					}
					copy(EBRNuevo.ERBName[:], "")
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
			if MBRTemporal.MRBPartitions[i].PART_Type[0] == 'e' {
				ParticionExtendida = &MBRTemporal.MRBPartitions[i]
				break
			}
		}
		if ParticionExtendida == nil {
			fmt.Println("Error: No se encontró una partición extendida para crear la partición lógica.")
			return
		}

		EBRPosterior := ParticionExtendida.PART_Start
		var EBRUltimo Estructura.EBR
		for {
			Utilidades.ReadObject(archivo, &EBRUltimo, int64(EBRPosterior))
			if strings.Contains(string(EBRUltimo.ERBName[:]), name) {
				fmt.Println("Error: El nombre: ", name, "ya está en uso en las particiones.")
				return
			}
			if EBRUltimo.ERBNext == -1 {
				break
			}
			EBRPosterior = EBRUltimo.ERBNext
		}

		var EBRNuevoPosterior int32
		if EBRUltimo.ERBSize == 0 {
			EBRNuevoPosterior = EBRPosterior
		} else {
			EBRNuevoPosterior = EBRUltimo.ERBStart + EBRUltimo.ERBSize
		}

		if EBRNuevoPosterior+int32(size)+int32(binary.Size(Estructura.EBR{})) > ParticionExtendida.PART_Start+ParticionExtendida.PART_Size {
			fmt.Println("Error: No hay suficiente espacio en la partición extendida para esta partición lógica.")
			return
		}

		if EBRUltimo.ERBSize != 0 {
			EBRUltimo.ERBNext = EBRNuevoPosterior
			Utilidades.WriteObject(archivo, EBRUltimo, int64(EBRPosterior))
		}

		newEBR := Estructura.EBR{
			ERBFit:   [1]byte{fit[0]},
			ERBStart: EBRNuevoPosterior + int32(binary.Size(Estructura.EBR{})),
			ERBSize:  int32(size),
			ERBNext:  -1,
		}
		copy(newEBR.ERBName[:], newEBR.ERBName[:])
		Utilidades.WriteObject(archivo, newEBR, int64(EBRNuevoPosterior))
		fmt.Println("Partición lógica creada exitosamente en la path: ", path, " con el nombre: ", name)
	}

	if err := Utilidades.WriteObject(archivo, MBRTemporal, 0); err != nil {
		return
	}
	defer archivo.Close()
	fmt.Println("======================================= FIN FDISK ======================================= ")
}

func Mount(ruta string, nombre string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "=======================================INICIO Mount=======================================")
	// ================================= Validar la ruta (path)
	if ruta == "" {
		fmt.Println("Error: La ruta es obligatoria.")
		return
	}
	// ================================= Validar el nombre (name)
	if nombre == "" {
		fmt.Println("Error: El nombre es obligatorio.")
		return
	}
	// ================================= Abrir archivo binario
	archivo, err := Utilidades.OpenFile(ruta)
	if err != nil {
		return
	}
	// ================================= Leer el MBR
	var MBRTemporal Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &MBRTemporal, 0); err != nil {
		return
	}
	// ================================= Verificar si la partición existe
	var ParticionExiste = false
	for i := 0; i < 4; i++ {
		if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Name[:]), nombre) {
			ParticionExiste = true
			break
		}
	}
	// ================================= Montar la partición
	if !ParticionExiste {
		fmt.Println("Error: No se encontró la partición con el nombre especificado.")
		return
	} else {
		numero, letra := NumeroLetraMountID(ruta)
		for i := 0; i < 4; i++ {
			if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Type[:]), "e") || strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Type[:]), "l") {
				fmt.Println("Error: No se puede montar una partición extendida o lógica.")
				return
			} else if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Name[:]), nombre) {
				MBRTemporal.MRBPartitions[i].PART_Status[0] = '1'
				MBRTemporal.MRBPartitions[i].PART_Id = [4]byte{'4', '7', byte(rune(letra)), byte('0' + numero)}
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
	fmt.Fprintln(buffer, "=================================MBR=================================")
	Estructura.PrintMBR(buffer, TempMBR2)
	fmt.Fprintln(buffer, "=================================MBR=================================")
	defer archivo.Close()
	fmt.Fprintln(buffer, "=======================================  End Mount =======================================  ")
}
