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

// Estructura para representar una partición montada
type PartitionMounted struct {
	Path     string
	Name     string
	ID       string
	Status   byte // 0: no montada, 1: montada
	LoggedIn bool // true: usuario ha iniciado sesión, false: no ha iniciado sesión
}

// Mapa para almacenar las particiones montadas, organizadas por disco
var mountedPartitions = make(map[string][]PartitionMounted)

// Función para imprimir las particiones montadas
func PrintMountedPartitions() {
	fmt.Println("Particiones montadas:")

	if len(mountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}

	for diskID, partitions := range mountedPartitions {
		fmt.Printf("Disco ID: %s\n", diskID)
		for _, partition := range partitions {
			loginStatus := "No"
			if partition.LoggedIn {
				loginStatus = "Sí"
			}
			fmt.Printf(" - Partición Name: %s, ID: %s, Path: %s, Status: %c, LoggedIn: %s\n",
				partition.Name, partition.ID, partition.Path, partition.Status, loginStatus)
		}
	}
	fmt.Println("")
}

// Función para obtener las particiones montadas
func GetMountedPartitions() map[string][]PartitionMounted {
	return mountedPartitions
}

// Función para marcar una partición como logueada
func MarkPartitionAsLoggedIn(id string) {
	for diskID, partitions := range mountedPartitions {
		for i, partition := range partitions {
			if partition.ID == id {
				mountedPartitions[diskID][i].LoggedIn = true
				fmt.Println("Partición con ID marcada como logueada.")
				return
			}
		}
	}
	fmt.Printf("No se encontró la partición con ID %s para marcarla como logueada.\n", id)
}

// Función para obtener el ID del último disco montado
func getLastDiskID() string {
	var lastDiskID string
	for diskID := range mountedPartitions {
		lastDiskID = diskID
	}
	return lastDiskID
}

func generateDiskID(path string) string {
	return strings.ToLower(path)
}

func Mkdisk(size int, fit string, unit string, path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "INICIO MKDISK=======================================")
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
	fmt.Fprintln(buffer, "Disco creado con éxito en la path: ", path)
	fmt.Fprintln(buffer, "End MKDISK=======================================")
}

func Rmdisk(path string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "INICIO RMDISK=======================================")

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
	fmt.Fprintln(buffer, "End RMDISK=======================================")
}

func Fdisk(size int, path string, name string, unit string, type_ string, fit string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "FDISK---------------------------------------------------------------------\n")
	// Validar el tamaño (size)
	if size <= 0 {
		fmt.Fprintf(buffer, "Error FDISK: EL tamaño de la partición debe ser mayor que 0.\n")
		return
	}
	// Validar la unit (unit)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Fprintf(buffer, "Error FDISK: La unit de tamaño debe ser Bytes, Kilobytes, Megabytes.\n")
		return
	}
	// Validar la path (path)
	if path == "" {
		fmt.Fprintf(buffer, "Error FDISK: La path del disco es obligatoria.\n")
		return
	}
	// Validar el type_ (type_)
	if type_ != "p" && type_ != "e" && type_ != "l" {
		fmt.Fprintf(buffer, "Error FDISK: El type_ de partición debe ser Primaria, Extendida, Lógica.\n")
		return
	}
	// Validar el fit (fit)
	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Fprintf(buffer, "Error FDISK: El fit de la partición debe ser BF, WF o FF.\n")
		return
	}
	// Validar el name (name)
	if name == "" {
		fmt.Fprintf(buffer, "Error FDISK: El name de la partición es obligatorio.\n")
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
		if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Name[:]), name) {
			fmt.Fprintf(buffer, "Error FDISK: El name: %s ya está en uso en las particiones.\n", name)
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

	if TotalParticiones >= 4 && type_ != "l" {
		fmt.Fprintf(buffer, "Error FDISK: No se pueden crear más de 4 particiones primarias o extendidas en total.\n")
		return
	}
	if type_ == "e" && ContadorExtendida > 0 {
		fmt.Fprintf(buffer, "Error FDISK: Solo se permite una partición extendida por disco.\n")
		return
	}
	if type_ == "l" && ContadorExtendida == 0 {
		fmt.Fprintf(buffer, "Error FDISK: No se puede crear una partición lógica sin una partición extendida.\n")
		return
	}
	if EspacioUtilizado+int32(size) > MBRTemporal.MRBSize {
		fmt.Fprintf(buffer, "Error FDISK: No hay suficiente espacio en el disco para crear esta partición.\n")
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
				copy(MBRTemporal.MRBPartitions[i].PART_Name[:], name)
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
					if err := Utilidades.WriteObject(archivo, EBRNuevo, int64(EBRInicio)); err != nil {
						return
					}
				}
				fmt.Fprintf(buffer, "Partición creada exitosamente en la path: %s con el name: %s.\n", path, name)
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
			fmt.Fprintf(buffer, "Error FDISK: No se encontró una partición extendida para crear la partición lógica.\n")
			return
		}

		EBRPosterior := ParticionExtendida.PART_Start
		var EBRUltimo Estructura.EBR
		for {
			if err := Utilidades.ReadObject(archivo, &EBRUltimo, int64(EBRPosterior)); err != nil {
				return
			}
			if strings.Contains(string(EBRUltimo.ERBName[:]), name) {
				fmt.Fprintf(buffer, "Error FDISK: El name: %s ya está en uso en las particiones.\n", name)
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
			fmt.Fprintf(buffer, "Error FDISK: No hay suficiente espacio en la partición extendida para esta partición lógica.\n")
			return
		}

		if EBRUltimo.ERBSize != 0 {
			EBRUltimo.ERBNext = EBRNuevoPosterior
			if err := Utilidades.WriteObject(archivo, EBRUltimo, int64(EBRPosterior)); err != nil {
				return
			}
		}

		newEBR := Estructura.EBR{
			ERBFit:   [1]byte{fit[0]},
			ERBStart: EBRNuevoPosterior + int32(binary.Size(Estructura.EBR{})),
			ERBSize:  int32(size),
			ERBNext:  -1,
		}
		copy(newEBR.ERBName[:], name)
		if err := Utilidades.WriteObject(archivo, newEBR, int64(EBRNuevoPosterior)); err != nil {
			return
		}
		fmt.Fprintf(buffer, "Partición lógica creada exitosamente en la path: %s con el name: %s.\n", path, name)
		fmt.Println("---------------------------------------------")
		EBRActual := ParticionExtendida.PART_Start
		for {
			var EBRTemp Estructura.EBR
			if err := Utilidades.ReadObject(archivo, &EBRTemp, int64(EBRActual)); err != nil {
				fmt.Fprintf(buffer, "Error leyendo EBR: %v\n", err)
				return
			}
			Estructura.PrintEBR(buffer, EBRTemp)
			if EBRTemp.ERBNext == -1 {
				break
			}
			EBRActual = EBRTemp.ERBNext
		}
		fmt.Println("---------------------------------------------")
	}
	if err := Utilidades.WriteObject(archivo, MBRTemporal, 0); err != nil {
		return
	}
	var TempMRB Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &TempMRB, 0); err != nil {
		return
	}
	fmt.Println("---------------------------------------------")
	Estructura.PrintMBR(buffer, TempMRB)
	fmt.Println("---------------------------------------------")
	defer archivo.Close()

	fmt.Println("FIN FDISK=======================================")
	fmt.Println("")
}

// Función para montar particiones
func Mount(path string, name string, buffer *bytes.Buffer) {
	fmt.Fprintln(buffer, "Start MOUNT=======================================")
	file, err := Utilidades.OpenFile(path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo en la path:", path)
		return
	}
	defer file.Close()

	var TempMBR Estructura.MRB
	if err := Utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR desde el archivo")
		return
	}

	fmt.Fprintf(buffer, "Buscando partición con name: '%s'\n", name)

	partitionFound := false
	var partition Estructura.Partition
	var partitionIndex int

	// Convertir el name a comparar a un arreglo de bytes de longitud fija
	nameBytes := [16]byte{}
	copy(nameBytes[:], []byte(name))

	for i := 0; i < 4; i++ {
		if TempMBR.MRBPartitions[i].PART_Type[0] == 'p' && bytes.Equal(TempMBR.MRBPartitions[i].PART_Name[:], nameBytes[:]) {
			partition = TempMBR.MRBPartitions[i]
			partitionIndex = i
			partitionFound = true
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: Partición no encontrada o no es una partición primaria")
		return
	}

	// Verificar si la partición ya está montada
	if partition.PART_Status[0] == '1' {
		fmt.Println("Error: La partición ya está montada")
		return
	}

	//fmt.Fprint("Partición encontrada: '%s' en posición %d\n", string(partition.Name[:]), partitionIndex+1)

	// Generar el ID de la partición
	diskID := generateDiskID(path)

	// Verificar si ya se ha montado alguna partición de este disco
	mountedPartitionsInDisk := mountedPartitions[diskID]
	var letter byte

	if len(mountedPartitionsInDisk) == 0 {
		// Es un nuevo disco, asignar la siguiente letra disponible
		if len(mountedPartitions) == 0 {
			letter = 'a'
		} else {
			lastDiskID := getLastDiskID()
			lastLetter := mountedPartitions[lastDiskID][0].ID[len(mountedPartitions[lastDiskID][0].ID)-1]
			letter = lastLetter + 1
		}
	} else {
		// Utilizar la misma letra que las otras particiones montadas en el mismo disco
		letter = mountedPartitionsInDisk[0].ID[len(mountedPartitionsInDisk[0].ID)-1]
	}

	// Incrementar el número para esta partición
	carnet := "202201947" // Cambiar su carnet aquí
	lastTwoDigits := carnet[len(carnet)-2:]
	partitionID := fmt.Sprintf("%s%d%c", lastTwoDigits, partitionIndex+1, letter)

	// Actualizar el estado de la partición a montada y asignar el ID
	partition.PART_Status[0] = '1'
	copy(partition.PART_Id[:], partitionID)
	TempMBR.MRBPartitions[partitionIndex] = partition
	mountedPartitions[diskID] = append(mountedPartitions[diskID], PartitionMounted{
		Path:   path,
		Name:   name,
		ID:     partitionID,
		Status: '1',
	})

	// Escribir el MBR actualizado al archivo
	if err := Utilidades.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo sobrescribir el MBR en el archivo")
		return
	}

	fmt.Fprintf(buffer, "Partición montada con ID: %s\n", partitionID)

	fmt.Println("")
	// Imprimir el MBR actualizado
	fmt.Println("MBR actualizado:")
	Estructura.PrintMBR(buffer, TempMBR)
	fmt.Println("")

	// Imprimir las particiones montadas (solo estan mientras dure la sesion de la consola)
	PrintMountedPartitions()
	fmt.Fprintln(buffer, "End MOUNT=======================================")
}
