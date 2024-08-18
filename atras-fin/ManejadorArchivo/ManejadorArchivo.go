package ManejadorArchivo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"proyecto1/Estructura"
	"proyecto1/Utilidades"
	"strings"
	"time"
)

func Mkfs(id string, type_ string, writer *bytes.Buffer) {
	fmt.Fprintln(writer, "=================================Start MKFS=================================")
	fmt.Fprintln(writer, "Id:", id)
	fmt.Fprintln(writer, "Type:", type_)

	archivo, err := Utilidades.OpenFile("/home/pablo03r/discosp1/disco1.mia")
	if err != nil {
		return
	}

	var MBRTemporal Estructura.MRB
	if err := Utilidades.ReadObject(archivo, &MBRTemporal, 0); err != nil {
		return
	}

	var IndiceParticion int = -1
	for i := 0; i < 4; i++ {
		if MBRTemporal.MRBPartitions[i].PART_Size != 0 {
			if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Id[:]), "47A1") {
				fmt.Println("La partición: ", id, " encontrada correctamente.")
				if strings.Contains(string(MBRTemporal.MRBPartitions[i].PART_Status[:]), "1") {
					fmt.Println("La partición: ", id, " está montada correctamente.")
					IndiceParticion = i
				} else {
					fmt.Println("La partición ", id, " no está montada correctamente.")
					return
				}
				break
			}
		}
	}

	if IndiceParticion == -1 {
		fmt.Println("La partición: ", id, " no existe.")
		return
	}

	numerador := int32(MBRTemporal.MRBPartitions[IndiceParticion].PART_Size - int32(binary.Size(Estructura.SuperBlock{})))
	denrominador_base := int32(4 + int32(binary.Size(Estructura.Inode{})) + 3*int32(binary.Size(Estructura.FileBlock{})))
	denrominador := denrominador_base
	n := int32(numerador / denrominador)

	var NuevoSuperBloque Estructura.SuperBlock
	NuevoSuperBloque.S_Inodes_Count = 0
	NuevoSuperBloque.S_Blocks_Count = 0
	NuevoSuperBloque.S_Free_Blocks_Count = 3 * n
	NuevoSuperBloque.S_Free_Inodes_Count = n
	FechaActual := time.Now()
	FechaCreacion := FechaActual.Format("2006-01-02")
	copy(NuevoSuperBloque.S_Mtime[:], FechaCreacion)
	copy(NuevoSuperBloque.S_Umtime[:], FechaCreacion)
	NuevoSuperBloque.S_Mnt_Count = 0
	fmt.Println("=================================End MKFS=================================")

	SistemaEXT2(n, MBRTemporal.MRBPartitions[IndiceParticion], NuevoSuperBloque, FechaCreacion, archivo, writer)
	defer archivo.Close()
}

func SistemaEXT2(n int32, Particion Estructura.Partition, NuevoSuperBloque Estructura.SuperBlock, Fecha string, archivo *os.File, writer io.Writer) {
	fmt.Println("=================================EXT2=================================")

	NuevoSuperBloque.S_Filesystem_Type = 2
	NuevoSuperBloque.S_BM_Inode_Start = Particion.PART_Start + int32(binary.Size(Estructura.SuperBlock{}))
	NuevoSuperBloque.S_BM_Block_Start = NuevoSuperBloque.S_BM_Inode_Start + n
	NuevoSuperBloque.S_Inode_Start = NuevoSuperBloque.S_BM_Block_Start + 3*n
	NuevoSuperBloque.S_Block_Start = NuevoSuperBloque.S_Inode_Start + n*int32(binary.Size(Estructura.Inode{}))
	// ================================= crear el super bloque =================================
	NuevoSuperBloque.S_Magic = 0xEF53
	NuevoSuperBloque.S_Mnt_Count = 1
	NuevoSuperBloque.S_Inode_Size = int32(binary.Size(Estructura.Inode{}))
	NuevoSuperBloque.S_Block_Size = int32(binary.Size(Estructura.FileBlock{}))
	//================================= Crear el bitmap de inodos y bloques =================================
	NuevoSuperBloque.S_Free_Inodes_Count -= 1
	NuevoSuperBloque.S_Free_Blocks_Count -= 1
	NuevoSuperBloque.S_Free_Inodes_Count -= 1
	NuevoSuperBloque.S_Free_Blocks_Count -= 1

	for i := int32(0); i < n; i++ {
		err := Utilidades.WriteObject(archivo, byte(0), int64(NuevoSuperBloque.S_BM_Inode_Start+i))
		if err != nil {
			return
		}
	}

	for i := int32(0); i < 3*n; i++ {
		err := Utilidades.WriteObject(archivo, byte(0), int64(NuevoSuperBloque.S_BM_Block_Start+i))
		if err != nil {
			return
		}
	}

	var NuevoInodo Estructura.Inode
	for i := int32(0); i < 15; i++ {
		NuevoInodo.I_Block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err := Utilidades.WriteObject(archivo, NuevoInodo, int64(NuevoSuperBloque.S_Inode_Start+i*int32(binary.Size(Estructura.Inode{}))))
		if err != nil {
			return
		}
	}

	var NuevoBloqueArchivo Estructura.FileBlock
	for i := int32(0); i < 3*n; i++ {
		err := Utilidades.WriteObject(archivo, NuevoBloqueArchivo, int64(NuevoSuperBloque.S_Block_Start+i*int32(binary.Size(Estructura.FileBlock{}))))
		if err != nil {
			return
		}
	}

	var Inodo0 Estructura.Inode
	//-------------
	//SE DEBE VER DE DONDE SALE EL USUARIO
	Inodo0.I_Uid = 1
	Inodo0.I_Gid = 0
	Inodo0.I_Size = int32(binary.Size(Estructura.Inode{}))
	//-------------

	copy(Inodo0.I_Atime[:], Fecha)
	copy(Inodo0.I_Ctime[:], Fecha)
	copy(Inodo0.I_Mtime[:], Fecha)
	Inodo0.I_Type = '1'
	copy(Inodo0.I_Perm[:], "664")
	for i := int32(0); i < 15; i++ {
		Inodo0.I_Block[i] = -1
	}
	Inodo0.I_Block[0] = 0

	//================================= Crear el bloque de carpetas =================================
	var BloqueCarpeta0 Estructura.FolderBlock
	copy(BloqueCarpeta0.B_Content[0].B_Name[:], ".")
	BloqueCarpeta0.B_Content[0].B_Inodo = 0
	copy(BloqueCarpeta0.B_Content[1].B_Name[:], "..")
	BloqueCarpeta0.B_Content[1].B_Inodo = 0
	copy(BloqueCarpeta0.B_Content[2].B_Name[:], "users.txt")
	BloqueCarpeta0.B_Content[2].B_Inodo = 1
	BloqueCarpeta0.B_Content[3].B_Inodo = -1
	// ==================================================================

	// ==================================================================
	NuevoSuperBloque.S_Inodes_Count++

	var Inodo1 Estructura.Inode //Inode 1
	Inodo1.I_Uid = 1
	Inodo1.I_Gid = 0
	Inodo1.I_Size = int32(binary.Size(Estructura.Inode{}))
	copy(Inodo1.I_Atime[:], Fecha)
	copy(Inodo1.I_Ctime[:], Fecha)
	copy(Inodo1.I_Mtime[:], Fecha)
	Inodo0.I_Type = '1'
	copy(Inodo1.I_Perm[:], "664")
	for i := int32(0); i < 15; i++ {
		Inodo1.I_Block[i] = -1
	}
	Inodo1.I_Block[0] = 1

	NuevoSuperBloque.S_Inodes_Count++

	DatosBloque := "1,G,root\n1,U,root,root,123\n"
	var BloqueArchivo1 Estructura.FileBlock
	copy(BloqueArchivo1.B_Content[:], DatosBloque)

	NuevoSuperBloque.S_First_Blo = int32(0)
	NuevoSuperBloque.S_First_Blo = int32(1)

	Utilidades.WriteObject(archivo, NuevoSuperBloque, int64(Particion.PART_Start))

	// ================================= Escritura de los bitmap de inodos =================================
	Utilidades.WriteObject(archivo, byte(1), int64(NuevoSuperBloque.S_BM_Inode_Start))
	Utilidades.WriteObject(archivo, byte(1), int64(NuevoSuperBloque.S_BM_Inode_Start+1))

	// ================================= Escritura de los bitmap de bloques =================================
	Utilidades.WriteObject(archivo, byte(1), int64(NuevoSuperBloque.S_BM_Block_Start))
	Utilidades.WriteObject(archivo, byte(1), int64(NuevoSuperBloque.S_BM_Block_Start+1))

	// ================================= Escritura de los inodos =================================
	Utilidades.WriteObject(archivo, Inodo0, int64(NuevoSuperBloque.S_Inode_Start))
	Utilidades.WriteObject(archivo, Inodo1, int64(NuevoSuperBloque.S_Inode_Start+int32(binary.Size(Estructura.Inode{}))))

	// ================================= Escritura de los bloques =================================
	Utilidades.WriteObject(archivo, BloqueCarpeta0, int64(NuevoSuperBloque.S_Block_Start))
	Utilidades.WriteObject(archivo, BloqueArchivo1, int64(NuevoSuperBloque.S_Block_Start+int32(binary.Size(Estructura.FileBlock{}))))

	fmt.Fprintln(writer, "=================================End EXT2=================================")
}
