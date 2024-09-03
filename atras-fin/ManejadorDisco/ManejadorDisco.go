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

// Función para imprimir las particiones montadas
func PrintMountedPartitions() {
	fmt.Println("Particiones montadas:")

	if len(mountedPartitions) == 0 {
		fmt.Println("No hay particiones montadas.")
		return
	}

	for diskID, partitions := range mountedPartitions {
		fmt.Println("Disco ID: %s\n", diskID)
		for _, partition := range partitions {
			loginStatus := "No"
			if partition.LoggedIn {
				loginStatus = "Sí"
			}
			fmt.Println(" - Partición Name: %s, ID: %s, Path: %s, Status: %c, LoggedIn: %s\n",
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

func Fdisk(size int, path string, name string, unit string, type_ string, fit string, buffer *bytes.Buffer) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Path:", path)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)

	// Validar fit (b/w/f)
	if fit != "bf" && fit != "ff" && fit != "wf" {
		fmt.Println("Error: Fit must be 'bf', 'ff', or 'wf'")
		return
	}

	// Validar size > 0
	if size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// Validar unit (b/k/m)
	if unit != "b" && unit != "k" && unit != "m" {
		fmt.Println("Error: Unit must be 'b', 'k', or 'm'")
		return
	}

	// Ajustar el tamaño en bytes
	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	// Abrir el archivo binario en la ruta proporcionada
	file, err := Utilidades.OpenFile(path)
	if err != nil {
		fmt.Println("Error: Could not open file at path:", path)
		return
	}

	var TempMBR Estructura.MRB
	// Leer el objeto desde el archivo binario
	if err := Utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file")
		return
	}

	// Imprimir el objeto MBR
	Estructura.PrintMBR(buffer, TempMBR)

	fmt.Println("-------------")

	// Validaciones de las particiones
	var primaryCount, extendedCount, totalPartitions int
	var usedSpace int32 = 0

	for i := 0; i < 4; i++ {
		if TempMBR.MRBPartitions[i].PART_Size != 0 {
			totalPartitions++
			usedSpace += TempMBR.MRBPartitions[i].PART_Size

			if TempMBR.MRBPartitions[i].PART_Type[0] == 'p' {
				primaryCount++
			} else if TempMBR.MRBPartitions[i].PART_Type[0] == 'e' {
				extendedCount++
			}
		}
	}

	// Validar que no se exceda el número máximo de particiones primarias y extendidas
	if totalPartitions >= 4 {
		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
		return
	}

	// Validar que solo haya una partición extendida
	if type_ == "e" && extendedCount > 0 {
		fmt.Println("Error: Solo se permite una partición extendida por disco.")
		return
	}

	// Validar que no se pueda crear una partición lógica sin una extendida
	if type_ == "l" && extendedCount == 0 {
		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
		return
	}

	// Validar que el tamaño de la nueva partición no exceda el tamaño del disco
	if usedSpace+int32(size) > TempMBR.MRBSize {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
		return
	}

	// Determinar la posición de inicio de la nueva partición
	var gap int32 = int32(binary.Size(TempMBR))
	if totalPartitions > 0 {
		gap = TempMBR.MRBPartitions[totalPartitions-1].PART_Start + TempMBR.MRBPartitions[totalPartitions-1].PART_Size
	}

	// Encontrar una posición vacía para la nueva partición
	for i := 0; i < 4; i++ {
		if TempMBR.MRBPartitions[i].PART_Size == 0 {
			if type_ == "p" || type_ == "e" {
				// Crear partición primaria o extendida
				TempMBR.MRBPartitions[i].PART_Size = int32(size)
				TempMBR.MRBPartitions[i].PART_Start = gap
				copy(TempMBR.MRBPartitions[i].PART_Name[:], name)
				copy(TempMBR.MRBPartitions[i].PART_Name[:], fit)
				copy(TempMBR.MRBPartitions[i].PART_Status[:], "0")
				copy(TempMBR.MRBPartitions[i].PART_Type[:], type_)
				TempMBR.MRBPartitions[i].PART_Correlative = int32(totalPartitions + 1)

				if type_ == "e" {
					// Inicializar el primer EBR en la partición extendida
					ebrStart := gap // El primer EBR se coloca al inicio de la partición extendida
					ebr := Estructura.EBR{
						ERBFit:   fit[0],
						ERBStart: ebrStart,
						ERBSize:  0,
						ERBNext:  -1,
					}
					copy(ebr.ERBName[:], "")
					Utilidades.WriteObject(file, ebr, int64(ebrStart))
				}

				break
			}
		}
	}

	// Manejar la creación de particiones lógicas dentro de una partición extendida
	if type_ == "l" {
		for i := 0; i < 4; i++ {
			if TempMBR.MRBPartitions[i].PART_Type[0] == 'e' {
				ebrPos := TempMBR.MRBPartitions[i].PART_Start
				var ebr Estructura.EBR
				for {
					Utilidades.ReadObject(file, &ebr, int64(ebrPos))
					if ebr.ERBNext == -1 {
						break
					}
					ebrPos = ebr.ERBNext
				}

				// Calcular la posición de inicio de la nueva partición lógica
				newEBRPos := ebr.ERBStart + ebr.ERBSize                      // El nuevo EBR se coloca después de la partición lógica anterior
				logicalPartitionStart := newEBRPos + int32(binary.Size(ebr)) // El inicio de la partición lógica es justo después del EBR

				// Ajustar el siguiente EBR
				ebr.ERBNext = newEBRPos
				Utilidades.WriteObject(file, ebr, int64(ebrPos))

				// Crear y escribir el nuevo EBR
				newEBR := Estructura.EBR{
					ERBFit:   fit[0],
					ERBStart: logicalPartitionStart,
					ERBSize:  int32(size),
					ERBNext:  -1,
				}
				copy(newEBR.ERBName[:], name)
				Utilidades.WriteObject(file, newEBR, int64(newEBRPos))

				// Imprimir el nuevo EBR creado
				fmt.Println("Nuevo EBR creado:")
				Estructura.PrintEBR(buffer, newEBR)
				fmt.Println("")

				// Imprimir todos los EBRs en la partición extendida
				fmt.Println("Imprimiendo todos los EBRs en la partición extendida:")
				ebrPos = TempMBR.MRBPartitions[i].PART_Start
				for {
					err := Utilidades.ReadObject(file, &ebr, int64(ebrPos))
					if err != nil {
						fmt.Println("Error al leer EBR:", err)
						break
					}
					Estructura.PrintEBR(buffer, ebr)
					if ebr.ERBNext == -1 {
						break
					}
					ebrPos = ebr.ERBNext
				}

				break
			}
		}
		fmt.Println("")
	}

	// Sobrescribir el MBR
	if err := Utilidades.WriteObject(file, TempMBR, 0); err != nil {
		fmt.Println("Error: Could not write MBR to file")
		return
	}

	var TempMBR2 Estructura.MRB
	// Leer el objeto nuevamente para verificar
	if err := Utilidades.ReadObject(file, &TempMBR2, 0); err != nil {
		fmt.Println("Error: Could not read MBR from file after writing")
		return
	}

	// Imprimir el objeto MBR actualizado
	Estructura.PrintMBR(buffer, TempMBR2)

	// Cerrar el archivo binario
	defer file.Close()

	fmt.Println("======FIN FDISK======")
	fmt.Println("")
}

// Función para montar particiones
func Mount(path string, name string, buffer *bytes.Buffer) {
	file, err := Utilidades.OpenFile(path)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo en la ruta:", path)
		return
	}
	defer file.Close()

	var TempMBR Estructura.MRB
	if err := Utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR desde el archivo")
		return
	}

	fmt.Fprint(buffer, "Buscando partición con nombre: '%s'\n", name)

	partitionFound := false
	var partition Estructura.Partition
	var partitionIndex int

	// Convertir el nombre a comparar a un arreglo de bytes de longitud fija
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

	fmt.Fprint(buffer, "Partición montada con ID: %s\n", partitionID)

	fmt.Println("")
	// Imprimir el MBR actualizado
	fmt.Println("MBR actualizado:")
	Estructura.PrintMBR(buffer, TempMBR)
	fmt.Println("")

	// Imprimir las particiones montadas (solo estan mientras dure la sesion de la consola)
	PrintMountedPartitions()
}

// func Fdisk(size int, unit string, path string, type_ string, fit string, name string, buffer *bytes.Buffer) {
// 	fmt.Println("======================================= Start FDISK =======================================")
// 	fmt.Fprintln(buffer, "Size:", size)
// 	fmt.Fprintln(buffer, "Path:", path)
// 	fmt.Fprintln(buffer, "Name:", name)
// 	fmt.Fprintln(buffer, "Unit:", unit)
// 	fmt.Fprintln(buffer, "Type:", type_)
// 	fmt.Fprintln(buffer, "Fit:", fit)

// 	// ================================= VALIDACIONES =================================
// 	if size <= 0 {
// 		fmt.Println("Error: Tamaño debe ser mayor que 0.")
// 		return
// 	}

// 	if unit != "b" && unit != "k" && unit != "m" {
// 		fmt.Println("Error: Unidad debe ser B, K, M.")
// 		return
// 	}

// 	if path == "" {
// 		fmt.Println("Error: La path es obligatoria.")
// 		return
// 	}

// 	if type_ != "p" && type_ != "e" && type_ != "l" {
// 		fmt.Println("Error: Tipo de partición debe ser P, E, L.")
// 		return
// 	}

// 	if fit != "bf" && fit != "wf" && fit != "ff" {
// 		fmt.Println("Error: Ajuste debe ser BF, WF o FF")
// 		return
// 	}

// 	if name == "" {
// 		fmt.Println("Error: El nombre es obligatorio.")
// 		return
// 	}

// 	if unit == "k" {
// 		size = size * 1024
// 	} else if unit == "m" {
// 		size = size * 1024 * 1024
// 	}

// 	// ================================= ABRIR ARCHIVO =================================
// 	archivo, err := Utilidades.OpenFile(path)
// 	if err != nil {
// 		return
// 	}

// 	// ================================= LEER MBR =================================
// 	var MBRTemporal Estructura.MRB
// 	if err := Utilidades.ReadObject(archivo, &MBRTemporal, 0); err != nil {
// 		return
// 	}

// 	// ================================= VALIDAR NOMBRE =================================
// 	for i := 0; i < 4; i++ {
// 		if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Name[:]), name) {
// 			fmt.Println("Error: El nombre: ", name, "ya está en uso en las particiones.")
// 			return
// 		}
// 	}

// 	var ContadorPrimaria, ContadorExtendida, TotalParticiones int
// 	var EspacioUtilizado int32 = 0

// 	for i := 0; i < 4; i++ {
// 		if MBRTemporal.MRBPartitions[i].PART_Size != 0 {
// 			TotalParticiones++
// 			EspacioUtilizado += MBRTemporal.MRBPartitions[i].PART_Size
// 			if MBRTemporal.MRBPartitions[i].PART_Type[0] == 'p' {
// 				ContadorPrimaria++
// 			} else if MBRTemporal.MRBPartitions[i].PART_Type[0] == 'e' {
// 				ContadorExtendida++
// 			}
// 		}
// 	}

// 	// ================================= VALIDAR PARTICIONES =================================
// 	if TotalParticiones >= 4 && type_ != "l" {
// 		fmt.Println("Error: No se pueden crear más de 4 particiones primarias o extendidas en total.")
// 		return
// 	}
// 	if type_ == "e" && ContadorExtendida > 0 {
// 		fmt.Println("Error: Solo se permite una partición extendida por disco.")
// 		return
// 	}
// 	if type_ == "l" && ContadorExtendida == 0 {
// 		fmt.Println("Error: No se puede crear una partición lógica sin una partición extendida.")
// 		return
// 	}
// 	if EspacioUtilizado+int32(size) > MBRTemporal.MRBSize {
// 		fmt.Println("Error: No hay suficiente espacio en el disco para crear esta partición.")
// 		return
// 	}

// 	var vacio int32 = int32(binary.Size(MBRTemporal))
// 	if TotalParticiones > 0 {
// 		vacio = MBRTemporal.MRBPartitions[TotalParticiones-1].PART_Start + MBRTemporal.MRBPartitions[TotalParticiones-1].PART_Size
// 	}

// 	for i := 0; i < 4; i++ {
// 		if MBRTemporal.MRBPartitions[i].PART_Size == 0 {
// 			if type_ == "p" || type_ == "e" {
// 				MBRTemporal.MRBPartitions[i].PART_Size = int32(size)
// 				MBRTemporal.MRBPartitions[i].PART_Start = vacio
// 				copy(MBRTemporal.MRBPartitions[i].PART_Name[:], []byte(name))
// 				copy(MBRTemporal.MRBPartitions[i].PART_Fit[:], fit)
// 				copy(MBRTemporal.MRBPartitions[i].PART_Status[:], "0")
// 				copy(MBRTemporal.MRBPartitions[i].PART_Type[:], type_)
// 				MBRTemporal.MRBPartitions[i].PART_Correlative = int32(TotalParticiones + 1)
// 				if type_ == "e" {
// 					EBRInicio := vacio
// 					EBRNuevo := Estructura.EBR{
// 						ERBFit:   [1]byte{fit[0]},
// 						ERBStart: EBRInicio,
// 						ERBSize:  0,
// 						ERBNext:  -1,
// 					}
// 					copy(EBRNuevo.ERBName[:], "")
// 					Utilidades.WriteObject(archivo, EBRNuevo, int64(EBRInicio))
// 				}
// 				fmt.Println(buffer, "Partición creada exitosamente en la path: ", path, " con el nombre: ", name)
// 				break
// 			}
// 		}
// 	}

// 	if type_ == "l" {
// 		var ParticionExtendida *Estructura.Partition
// 		for i := 0; i < 4; i++ {
// 			if MBRTemporal.MRBPartitions[i].PART_Type[0] == 'e' {
// 				ParticionExtendida = &MBRTemporal.MRBPartitions[i]
// 				break
// 			}
// 		}
// 		if ParticionExtendida == nil {
// 			fmt.Println("Error: No se encontró una partición extendida para crear la partición lógica.")
// 			return
// 		}

// 		EBRPosterior := ParticionExtendida.PART_Start
// 		var EBRUltimo Estructura.EBR
// 		for {
// 			Utilidades.ReadObject(archivo, &EBRUltimo, int64(EBRPosterior))
// 			if strings.Contains(string(EBRUltimo.ERBName[:]), name) {
// 				fmt.Println("Error: El nombre: ", name, "ya está en uso en las particiones.")
// 				return
// 			}
// 			if EBRUltimo.ERBNext == -1 {
// 				break
// 			}
// 			EBRPosterior = EBRUltimo.ERBNext
// 		}

// 		var EBRNuevoPosterior int32
// 		if EBRUltimo.ERBSize == 0 {
// 			EBRNuevoPosterior = EBRPosterior
// 		} else {
// 			EBRNuevoPosterior = EBRUltimo.ERBStart + EBRUltimo.ERBSize
// 		}

// 		if EBRNuevoPosterior+int32(size)+int32(binary.Size(Estructura.EBR{})) > ParticionExtendida.PART_Start+ParticionExtendida.PART_Size {
// 			fmt.Println("Error: No hay suficiente espacio en la partición extendida para esta partición lógica.")
// 			return
// 		}

// 		if EBRUltimo.ERBSize != 0 {
// 			EBRUltimo.ERBNext = EBRNuevoPosterior
// 			Utilidades.WriteObject(archivo, EBRUltimo, int64(EBRPosterior))
// 		}

// 		newEBR := Estructura.EBR{
// 			ERBFit:   [1]byte{fit[0]},
// 			ERBStart: EBRNuevoPosterior + int32(binary.Size(Estructura.EBR{})),
// 			ERBSize:  int32(size),
// 			ERBNext:  -1,
// 		}
// 		copy(newEBR.ERBName[:], newEBR.ERBName[:])
// 		Utilidades.WriteObject(archivo, newEBR, int64(EBRNuevoPosterior))
// 		fmt.Println("Partición lógica creada exitosamente en la path: ", path, " con el nombre: ", name)
// 	}

// 	if err := Utilidades.WriteObject(archivo, MBRTemporal, 0); err != nil {
// 		return
// 	}
// 	defer archivo.Close()
// 	fmt.Println("======================================= FIN FDISK ======================================= ")
// }

// func Mount(ruta string, nombre string, buffer *bytes.Buffer) {
// 	fmt.Fprintln(buffer, "=======================================INICIO Mount=======================================")
// 	// ================================= Validar la ruta (path)
// 	if ruta == "" {
// 		fmt.Println("Error: La ruta es obligatoria.")
// 		return
// 	}
// 	// ================================= Validar el nombre (name)
// 	if nombre == "" {
// 		fmt.Println("Error: El nombre es obligatorio.")
// 		return
// 	}
// 	// ================================= Abrir archivo binario
// 	archivo, err := Utilidades.OpenFile(ruta)
// 	if err != nil {
// 		return
// 	}
// 	// ================================= Leer el MBR
// 	var MBRTemporal Estructura.MRB
// 	if err := Utilidades.ReadObject(archivo, &MBRTemporal, 0); err != nil {
// 		return
// 	}
// 	// ================================= Verificar si la partición existe
// 	var ParticionExiste = false
// 	for i := 0; i < 4; i++ {
// 		if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Name[:]), nombre) {
// 			ParticionExiste = true
// 			break
// 		}
// 	}
// 	// ================================= Montar la partición
// 	if !ParticionExiste {
// 		fmt.Println("Error: No se encontró la partición con el nombre especificado.")
// 		return
// 	} else {
// 		numero, letra := NumeroLetraMountID(ruta)
// 		for i := 0; i < 4; i++ {
// 			if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Type[:]), "e") || strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Type[:]), "l") {
// 				fmt.Println("Error: No se puede montar una partición extendida o lógica.")
// 				return
// 			} else if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Name[:]), nombre) {
// 				MBRTemporal.MRBPartitions[i].PART_Status[0] = '1'
// 				MBRTemporal.MRBPartitions[i].PART_Id = [4]byte{'4', '7', byte(rune(letra)), byte('0' + numero)}
// 				fmt.Println("Partición montada exitosamente en la ruta: ", ruta)
// 				break
// 			}
// 		}
// 		IncrementoNumeroMountID(ruta)
// 	}
// 	if err := Utilidades.WriteObject(archivo, MBRTemporal, 0); err != nil {
// 		return
// 	}
// 	var TempMBR2 Estructura.MRB
// 	if err := Utilidades.ReadObject(archivo, &TempMBR2, 0); err != nil {
// 		return
// 	}
// 	fmt.Fprintln(buffer, "=================================MBR=================================")
// 	Estructura.PrintMBR(buffer, TempMBR2)
// 	fmt.Fprintln(buffer, "=================================MBR=================================")
// 	defer archivo.Close()
// 	fmt.Fprintln(buffer, "=======================================  End Mount =======================================  ")
// }
