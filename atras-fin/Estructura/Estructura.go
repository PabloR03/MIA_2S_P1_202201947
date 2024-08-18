package Estructura

import (
	"bytes"
	"fmt"
)

//  =================================Estructura MRB=================================

type MRB struct {
	MRBSize         int32
	MRBCreationDate [10]byte
	MRBSignature    int32
	MRBFit          [1]byte
	MRBPartitions   [4]Partition
}

func PrintMBR(buffer *bytes.Buffer, data MRB) {
	fmt.Fprintf(buffer, "Fecha de Creación: %s, Ajuste: %s, Tamaño: %d, Identificador: %d\n",
		string(data.MRBCreationDate[:]), string(data.MRBFit[:]), data.MRBSize, data.MRBSignature)
	for i := 0; i < 4; i++ {
		PrintPartition(buffer, data.MRBPartitions[i])
	}
}

//  =================================Estructura Particion=================================

type Partition struct {
	PART_Status      [1]byte
	PART_Type        [1]byte
	PART_Fit         [1]byte
	PART_Start       int32
	PART_Size        int32
	PART_Name        [16]byte
	PART_Correlative int32
	PART_Id          [4]byte
	PART_Unit        [1]byte
	PART_Path        [100]byte
}

func PrintPartition(buffer *bytes.Buffer, data Partition) {
	fmt.Fprintf(buffer, "Nombre: %s, Tipo: %s, Inicio: %d, Tamaño: %d, Estado: %s, ID: %s, Ajuste: %s, Correlativo: %d\n",
		string(data.PART_Name[:]), string(data.PART_Type[:]), data.PART_Start, data.PART_Size, string(data.PART_Status[:]),
		string(data.PART_Id[:]), string(data.PART_Fit[:]), data.PART_Correlative)
}

//  =================================Estructura EBR=================================

type EBR struct {
	ERBMount [1]byte
	ERBFit   [1]byte
	ERBStart int32
	ERBSize  int32
	ERBNext  int32
	ERBName  [16]byte
}

func PrintEBR(buffer *bytes.Buffer, data EBR) {
	fmt.Fprintf(buffer, "Mount: %s, Fit: %s, Start: %d, Size: %d, Next: %d, Name: %s\n",
		string(data.ERBMount[:]), string(data.ERBFit[:]), data.ERBStart, data.ERBSize, data.ERBNext, string(data.ERBName[:]))
}

// =================================Estuctura MountId=================================

type MountId struct {
	MIDpath   string
	MIDnumber int32
	MIDletter int32
}

// =================================Estuctura Superblock=================================

type SuperBlock struct {
	S_Filesystem_Type   int32
	S_Inodes_Count      int32
	S_Blocks_Count      int32
	S_Free_Blocks_Count int32
	S_Free_Inodes_Count int32
	S_Mtime             [17]byte
	S_Umtime            [17]byte
	S_Mnt_Count         int32
	S_Magic             int32
	S_Inode_Size        int32
	S_Block_Size        int32
	S_Fist_Ino          int32
	S_First_Blo         int32
	S_BM_Inode_Start    int32
	S_BM_Block_Start    int32
	S_Inode_Start       int32
	S_Block_Start       int32
}

// =================================Estuctura Inode=================================

type Inode struct {
	I_Uid   int32
	I_Gid   int32
	I_Size  int32
	I_Atime [17]byte
	I_Ctime [17]byte
	I_Mtime [17]byte
	I_Block [15]int32
	I_Type  byte
	I_Perm  [3]byte
}

func PrintInode(buffer *bytes.Buffer, data Inode) {
	fmt.Fprintf(buffer, "INODO %d\nUID: %d \nGID: %d \nSIZE: %d \nACTUAL DATE: %s \nCREATION TIME: %s \nMODIFY TIME: %s \nBLOCKS:%d \nTYPE:%s \nPERM:%s \n",
		int(data.I_Gid),
		int(data.I_Uid),
		int(data.I_Gid),
		int(data.I_Size),
		data.I_Atime[:],
		data.I_Ctime[:],
		data.I_Mtime[:],
		data.I_Block[:],
		string(data.I_Type),
		string(data.I_Perm[:]),
	)
}

// =================================Estuctura Fileblock=================================

type FileBlock struct {
	B_Content [64]byte
}

// =================================Estuctura Folderblock=================================

type FolderBlock struct {
	B_Content [4]Content
}

type Content struct {
	B_Name  [12]byte
	B_Inodo int32
}

type PointerBlock struct {
	B_Pointers [16]int32
}
