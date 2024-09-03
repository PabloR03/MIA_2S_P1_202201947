package Usuario

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"proyecto1/Estructura"
	"proyecto1/ManejadorDisco"
	"proyecto1/Utilidades"
	"strings"
)

func Login(user string, pass string, id string, writer *bytes.Buffer) {
	fmt.Println("======Start LOGIN======")
	fmt.Println("User:", user)
	fmt.Println("Pass:", pass)
	fmt.Println("Id:", id)

	// Verificar si el usuario ya está logueado buscando en las particiones montadas
	mountedPartitions := ManejadorDisco.GetMountedPartitions()
	var filepath string
	var partitionFound bool
	var login bool = false

	for _, partitions := range mountedPartitions {
		for _, partition := range partitions {
			if partition.ID == id && partition.LoggedIn { // Verifica si ya está logueado
				fmt.Println("Ya existe un usuario logueado!")
				return
			}
			if partition.ID == id { // Encuentra la partición correcta
				filepath = partition.Path
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Println("Error: No se encontró ninguna partición montada con el ID proporcionado")
		return
	}

	// Abrir archivo binario
	file, err := Utilidades.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error: No se pudo abrir el archivo:", err)
		return
	}
	defer file.Close()

	var TempMBR Estructura.MRB
	// Leer el MBR desde el archivo binario
	if err := Utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error: No se pudo leer el MBR:", err)
		return
	}

	// Imprimir el MBR
	Estructura.PrintMBR(writer, TempMBR)
	fmt.Println("-------------")

	var index int = -1
	// Iterar sobre las particiones del MBR para encontrar la correcta
	for i := 0; i < 4; i++ {
		if TempMBR.MRBPartitions[i].PART_Size != 0 {
			if strings.Contains(string(TempMBR.MRBPartitions[i].PART_Id[:]), id) {
				fmt.Println("Partition found")
				if TempMBR.MRBPartitions[i].PART_Status[0] == '1' {
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		Estructura.PrintPartition(writer, TempMBR.MRBPartitions[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	var tempSuperblock Estructura.SuperBlock
	// Leer el Superblock desde el archivo binario
	if err := Utilidades.ReadObject(file, &tempSuperblock, int64(TempMBR.MRBPartitions[index].PART_Start)); err != nil {
		fmt.Println("Error: No se pudo leer el Superblock:", err)
		return
	}

	// Buscar el archivo de usuarios /users.txt -> retorna índice del Inodo
	indexInode := InitSearch("/users.txt", file, tempSuperblock)

	var crrInode Estructura.Inode
	// Leer el Inodo desde el archivo binario
	if err := Utilidades.ReadObject(file, &crrInode, int64(tempSuperblock.S_Inode_Start+indexInode*int32(binary.Size(Estructura.Inode{})))); err != nil {
		fmt.Println("Error: No se pudo leer el Inodo:", err)
		return
	}

	// Leer datos del archivo
	data := GetInodeFileData(crrInode, file, tempSuperblock)

	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// Iterar a través de las líneas para verificar las credenciales
	for _, line := range lines {
		words := strings.Split(line, ",")
		if len(words) == 5 {
			if (strings.Contains(words[3], user)) && (strings.Contains(words[4], pass)) {
				login = true
				break
			}
		}
	}

	// Imprimir información del Inodo
	fmt.Println("Inode", crrInode.I_Block)

	// Si las credenciales son correctas y marcamos como logueado
	if login {
		fmt.Println("Usuario logueado con exito")
		ManejadorDisco.MarkPartitionAsLoggedIn(id) // Marcar la partición como logueada
	}

	fmt.Println("======End LOGIN======")
}

func InitSearch(path string, file *os.File, tempSuperblock Estructura.SuperBlock) int32 {
	fmt.Println("======Start BUSQUEDA INICIAL ======")
	fmt.Println("path:", path)
	// path = "/ruta/nueva"

	// split the path by /
	TempStepsPath := strings.Split(path, "/")
	StepsPath := TempStepsPath[1:]

	fmt.Println("StepsPath:", StepsPath, "len(StepsPath):", len(StepsPath))
	for _, step := range StepsPath {
		fmt.Println("step:", step)
	}

	var Inode0 Estructura.Inode
	// Read object from bin file
	if err := Utilidades.ReadObject(file, &Inode0, int64(tempSuperblock.S_Inode_Start)); err != nil {
		return -1
	}

	fmt.Println("======End BUSQUEDA INICIAL======")

	return SarchInodeByPath(StepsPath, Inode0, file, tempSuperblock)
}

// stack
func pop(s *[]string) string {
	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last
}

func SarchInodeByPath(StepsPath []string, Inode Estructura.Inode, file *os.File, tempSuperblock Estructura.SuperBlock) int32 {
	fmt.Println("======Start BUSQUEDA INODO POR PATH======")
	index := int32(0)
	SearchedName := strings.Replace(pop(&StepsPath), " ", "", -1)

	fmt.Println("========== SearchedName:", SearchedName)

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_Block {
		if block != -1 {
			if index < 13 {
				//CASO DIRECTO

				var crrFolderBlock Estructura.FolderBlock
				// Read object from bin file
				if err := Utilidades.ReadObject(file, &crrFolderBlock, int64(tempSuperblock.S_Block_Start+block*int32(binary.Size(Estructura.FolderBlock{})))); err != nil {
					return -1
				}

				for _, folder := range crrFolderBlock.B_Content {
					// fmt.Println("Folder found======")
					fmt.Println("Folder === Name:", string(folder.B_Name[:]), "B_inodo", folder.B_Inodo)

					if strings.Contains(string(folder.B_Name[:]), SearchedName) {

						fmt.Println("len(StepsPath)", len(StepsPath), "StepsPath", StepsPath)
						if len(StepsPath) == 0 {
							fmt.Println("Folder found======")
							return folder.B_Inodo
						} else {
							fmt.Println("NextInode======")
							var NextInode Estructura.Inode
							// Read object from bin file
							if err := Utilidades.ReadObject(file, &NextInode, int64(tempSuperblock.S_Inode_Start+folder.B_Inodo*int32(binary.Size(Estructura.Inode{})))); err != nil {
								return -1
							}
							return SarchInodeByPath(StepsPath, NextInode, file, tempSuperblock)
						}
					}
				}

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End BUSQUEDA INODO POR PATH======")
	return 0
}

func GetInodeFileData(Inode Estructura.Inode, file *os.File, tempSuperblock Estructura.SuperBlock) string {
	fmt.Println("======Start CONTENIDO DEL BLOQUE======")
	index := int32(0)
	// define content as a string
	var content string

	// Iterate over i_blocks from Inode
	for _, block := range Inode.I_Block {
		if block != -1 {
			//Dentro de los directos
			if index < 13 {
				var crrFileBlock Estructura.FileBlock
				// Read object from bin file
				if err := Utilidades.ReadObject(file, &crrFileBlock, int64(tempSuperblock.S_Block_Start+block*int32(binary.Size(Estructura.FileBlock{})))); err != nil {
					return ""
				}

				content += string(crrFileBlock.B_Content[:])

			} else {
				fmt.Print("indirectos")
			}
		}
		index++
	}

	fmt.Println("======End CONTENIDO DEL BLOQUE======")
	return content
}

// MKUSER
func AppendToFileBlock(inode *Estructura.Inode, newData string, file *os.File, superblock Estructura.SuperBlock) error {
	// Leer el contenido existente del archivo utilizando la función GetInodeFileData
	existingData := GetInodeFileData(*inode, file, superblock)

	// Concatenar el nuevo contenido
	fullData := existingData + newData

	// Asegurarse de que el contenido no exceda el tamaño del bloque
	if len(fullData) > len(inode.I_Block)*binary.Size(Estructura.FileBlock{}) {
		// Si el contenido excede, necesitas manejar bloques adicionales
		return fmt.Errorf("el tamaño del archivo excede la capacidad del bloque actual y no se ha implementado la creación de bloques adicionales")
	}

	// Escribir el contenido actualizado en el bloque existente
	var updatedFileBlock Estructura.FileBlock
	copy(updatedFileBlock.B_Content[:], fullData)
	if err := Utilidades.WriteObject(file, updatedFileBlock, int64(superblock.S_Block_Start+inode.I_Block[0]*int32(binary.Size(Estructura.FileBlock{})))); err != nil {
		return fmt.Errorf("error al escribir el bloque actualizado: %v", err)
	}

	// Actualizar el tamaño del inodo
	inode.I_Size = int32(len(fullData))
	if err := Utilidades.WriteObject(file, *inode, int64(superblock.S_Inode_Start+inode.I_Block[0]*int32(binary.Size(Estructura.Inode{})))); err != nil {
		return fmt.Errorf("error al actualizar el inodo: %v", err)
	}

	return nil
}
