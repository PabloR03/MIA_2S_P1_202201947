package ManejadorArchivo

import (
	"bytes"
	"encoding/binary"
	"fmt"

	//"io"
	"os"
	"proyecto1/Estructura"
	"proyecto1/ManejadorDisco"
	"proyecto1/Utilidades"
	"strings"
	"time"
)

// YA REVISADO
func Mkfs(id string, type_ string, fs_ string, writer *bytes.Buffer) {
	fmt.Fprintln(writer, "=-=-=-=-=-=-=-= INCIO MKFS =-=-=-=-=-=-=-=-=")
	fmt.Fprintln(writer, "Id:", id)
	fmt.Fprintln(writer, "Type:", type_)
	fmt.Fprintln(writer, "Fs:", fs_)
	println("Id:", id)
	println("Type:", type_)
	println("Fs:", fs_)

	// Buscar la partición montada por ID
	var mountedPartition ManejadorDisco.PartitionMounted
	var partitionFound bool

	for _, partitions := range ManejadorDisco.GetMountedPartitions() {
		for _, partition := range partitions {
			if partition.ID == id {
				mountedPartition = partition
				partitionFound = true
				break
			}
		}
		if partitionFound {
			break
		}
	}

	if !partitionFound {
		fmt.Fprintln(writer, "Particion no encontrada")
		return
	}

	if mountedPartition.Status != '1' { // Verifica si la partición está montada
		fmt.Fprintln(writer, "La particion aun no esta montada")
		return
	}

	// Abrir archivo binario
	file, err := Utilidades.OpenFile(mountedPartition.Path)
	if err != nil {
		return
	}

	var TempMBR Estructura.MRB
	// Leer objeto desde archivo binario
	if err := Utilidades.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Imprimir objeto
	Estructura.PrintMBR(writer, TempMBR)
	Estructura.PrintMBRnormal(TempMBR)

	fmt.Println("-------------")

	var index int = -1
	// Iterar sobre las particiones para encontrar la que tiene el nombre correspondiente
	for i := 0; i < 4; i++ {
		if TempMBR.MRBPartitions[i].PART_Size != 0 {
			if strings.Contains(string(TempMBR.MRBPartitions[i].PART_Id[:]), id) {
				index = i
				break
			}
		}
	}

	if index != -1 {
		Estructura.PrintPartition(writer, TempMBR.MRBPartitions[index])

	} else {
		fmt.Fprintln(writer, "Particion no encontrada (2)")
		return
	}

	numerador := int32(TempMBR.MRBPartitions[index].PART_Size - int32(binary.Size(Estructura.SuperBlock{})))
	denominador_base := int32(4 + int32(binary.Size(Estructura.Inode{})) + 3*int32(binary.Size(Estructura.FileBlock{})))
	var temp int32 = 0
	if fs_ == "2fs" {
		temp = 0
	} else {
		fmt.Fprintln(writer, "Error por el momento solo está disponible 2FS.")
	}
	denominador := denominador_base + temp
	n := int32(numerador / denominador)

	fmt.Println("INODOS:", n)

	// Crear el Superblock con todos los campos calculados
	var newSuperblock Estructura.SuperBlock
	newSuperblock.S_Filesystem_Type = 2 // EXT2
	newSuperblock.S_Inodes_Count = n
	newSuperblock.S_Blocks_Count = 3 * n
	newSuperblock.S_Free_Blocks_Count = 3*n - 2
	newSuperblock.S_Free_Inodes_Count = n - 2
	FechaActual := time.Now()
	FechaString := FechaActual.Format("02-01-2006 15:04:05")
	FechaBytes := []byte(FechaString)
	copy(newSuperblock.S_Mtime[:], FechaBytes)
	copy(newSuperblock.S_Umtime[:], FechaBytes)
	newSuperblock.S_Mnt_Count = 1
	newSuperblock.S_Magic = 0xEF53
	newSuperblock.S_Inode_Size = int32(binary.Size(Estructura.Inode{}))
	newSuperblock.S_Block_Size = int32(binary.Size(Estructura.FileBlock{}))

	// Calcula las posiciones de inicio
	newSuperblock.S_BM_Inode_Start = TempMBR.MRBPartitions[index].PART_Start + int32(binary.Size(Estructura.SuperBlock{}))
	newSuperblock.S_BM_Block_Start = newSuperblock.S_BM_Inode_Start + n
	newSuperblock.S_Inode_Start = newSuperblock.S_BM_Block_Start + 3*n
	newSuperblock.S_Block_Start = newSuperblock.S_Inode_Start + n*newSuperblock.S_Inode_Size

	if fs_ == "2fs" {
		create_ext2(n, TempMBR.MRBPartitions[index], newSuperblock, string(FechaBytes), file, writer)
	} else {
		fmt.Fprintln(writer, "EXT3 no está soportado.")
	}

	// Cerrar archivo binario
	defer file.Close()

	fmt.Fprintln(writer, "=-=-=-=-=-=-=-= FIN MKFS=-=-=-=-=-=-=-=-=")
}

func create_ext2(n int32, partition Estructura.Partition, newSuperblock Estructura.SuperBlock, date string, file *os.File, writer *bytes.Buffer) {
	fmt.Println("======Start CREATE EXT2======")
	fmt.Println("INODOS:", n)

	// Imprimir Superblock inicial
	Estructura.PrintSuperBlock(writer, newSuperblock)
	fmt.Println("Date:", date)

	// Escribe los bitmaps de inodos y bloques en el archivo
	for i := int32(0); i < n; i++ {
		if err := Utilidades.WriteObject(file, byte(0), int64(newSuperblock.S_BM_Inode_Start+i)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	for i := int32(0); i < 3*n; i++ {
		if err := Utilidades.WriteObject(file, byte(0), int64(newSuperblock.S_BM_Block_Start+i)); err != nil {
			fmt.Fprint(writer, "Error: ", err)
			return
		}
	}

	// Inicializa inodos y bloques con valores predeterminados
	if err := initInodesAndBlocks(n, newSuperblock, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Crea la carpeta raíz y el archivo users.txt
	if err := createRootAndUsersFile(newSuperblock, date, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Escribe el superbloque actualizado al archivo
	if err := Utilidades.WriteObject(file, newSuperblock, int64(partition.PART_Start)); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Marca los primeros inodos y bloques como usados
	if err := markUsedInodesAndBlocks(newSuperblock, file); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Leer e imprimir los inodos después de formatear
	fmt.Fprint(writer, "====== Imprimiendo Inodos ======")
	for i := int32(0); i < n; i++ {
		var inode Estructura.Inode
		offset := int64(newSuperblock.S_Inode_Start + i*int32(binary.Size(Estructura.Inode{})))
		if err := Utilidades.ReadObject(file, &inode, offset); err != nil {
			fmt.Println("Error al leer inodo: ", err)
			return
		}
		Estructura.PrintInode(writer, inode)
	}

	// Leer e imprimir los Folderblocks y Fileblocks después de formatear
	fmt.Println("====== Imprimiendo Folderblocks y Fileblocks ======")

	// Imprimir Folderblocks
	for i := int32(0); i < 1; i++ {
		var folderblock Estructura.FolderBlock
		offset := int64(newSuperblock.S_Block_Start + i*int32(binary.Size(Estructura.FolderBlock{})))
		if err := Utilidades.ReadObject(file, &folderblock, offset); err != nil {
			fmt.Println("Error al leer Folderblock: ", err)
			return
		}
		Estructura.PrintFolderBlock(writer, folderblock)
	}

	// Imprimir Fileblocks
	for i := int32(0); i < 1; i++ {
		var fileblock Estructura.FileBlock
		offset := int64(newSuperblock.S_Block_Start + int32(binary.Size(Estructura.FolderBlock{})) + i*int32(binary.Size(Estructura.FileBlock{})))
		if err := Utilidades.ReadObject(file, &fileblock, offset); err != nil {
			fmt.Println("Error al leer Fileblock: ", err)
			return
		}
		Estructura.PrintFileBlock(writer, fileblock)
	}

	// Imprimir el Superblock final
	Estructura.PrintSuperBlock(writer, newSuperblock)

	fmt.Fprint(writer, "======End CREATE EXT2======")
}

// Función auxiliar para inicializar inodos y bloques
func initInodesAndBlocks(n int32, newSuperblock Estructura.SuperBlock, file *os.File) error {
	var newInode Estructura.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_Block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		if err := Utilidades.WriteObject(file, newInode, int64(newSuperblock.S_Inode_Start+i*int32(binary.Size(Estructura.Inode{})))); err != nil {
			return err
		}
	}

	var newFileblock Estructura.FileBlock
	for i := int32(0); i < 3*n; i++ {
		if err := Utilidades.WriteObject(file, newFileblock, int64(newSuperblock.S_Block_Start+i*int32(binary.Size(Estructura.FileBlock{})))); err != nil {
			return err
		}
	}

	return nil
}

// Función auxiliar para crear la carpeta raíz y el archivo users.txt
func createRootAndUsersFile(newSuperblock Estructura.SuperBlock, date string, file *os.File) error {
	var Inode0, Inode1 Estructura.Inode
	initInode(&Inode0, date)
	initInode(&Inode1, date)

	Inode0.I_Block[0] = 0
	Inode1.I_Block[0] = 1

	// Asignar el tamaño real del contenido
	data := "1,G,root\n1,U,root,root,123\n"
	actualSize := int32(len(data))
	Inode1.I_Size = actualSize // Esto ahora refleja el tamaño real del contenido

	var Fileblock1 Estructura.FileBlock
	copy(Fileblock1.B_Content[:], data) // Copia segura de datos a Fileblock

	var Folderblock0 Estructura.FolderBlock
	Folderblock0.B_Content[0].B_Inodo = 0
	copy(Folderblock0.B_Content[0].B_Name[:], ".")
	Folderblock0.B_Content[1].B_Inodo = 0
	copy(Folderblock0.B_Content[1].B_Name[:], "..")
	Folderblock0.B_Content[2].B_Inodo = 1
	copy(Folderblock0.B_Content[2].B_Name[:], "users.txt")

	// Escribir los inodos y bloques en las posiciones correctas
	if err := Utilidades.WriteObject(file, Inode0, int64(newSuperblock.S_Inode_Start)); err != nil {
		return err
	}
	if err := Utilidades.WriteObject(file, Inode1, int64(newSuperblock.S_Inode_Start+int32(binary.Size(Estructura.Inode{})))); err != nil {
		return err
	}
	if err := Utilidades.WriteObject(file, Folderblock0, int64(newSuperblock.S_Block_Start)); err != nil {
		return err
	}
	if err := Utilidades.WriteObject(file, Fileblock1, int64(newSuperblock.S_Block_Start+int32(binary.Size(Estructura.FolderBlock{})))); err != nil {
		return err
	}

	return nil
}

// Función auxiliar para inicializar un inodo
func initInode(inode *Estructura.Inode, date string) {
	inode.I_Uid = 1
	inode.I_Gid = 1
	inode.I_Size = 0
	copy(inode.I_Atime[:], date)
	copy(inode.I_Ctime[:], date)
	copy(inode.I_Mtime[:], date)
	copy(inode.I_Perm[:], "664")

	for i := int32(0); i < 15; i++ {
		inode.I_Block[i] = -1
	}
}

// Función auxiliar para marcar los inodos y bloques usados
func markUsedInodesAndBlocks(newSuperblock Estructura.SuperBlock, file *os.File) error {
	if err := Utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_BM_Inode_Start)); err != nil {
		return err
	}
	if err := Utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_BM_Inode_Start+1)); err != nil {
		return err
	}
	if err := Utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_BM_Block_Start)); err != nil {
		return err
	}
	if err := Utilidades.WriteObject(file, byte(1), int64(newSuperblock.S_BM_Block_Start+1)); err != nil {
		return err
	}
	return nil
}
