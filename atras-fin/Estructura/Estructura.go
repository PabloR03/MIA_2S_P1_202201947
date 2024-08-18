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

func PrintMBR(data MRB, buffer *bytes.Buffer) {
	fmt.Printf("Fecha de Creación: %s, Ajuste: %s, Tamaño: %d, Identificador: %d\n",
		string(data.MRBCreationDate[:]), string(data.MRBFit[:]), data.MRBSize, data.MRBSignature)
	for i := 0; i < 4; i++ {
		PrintPartition(data.MRBPartitions[i], buffer)
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

func PrintPartition(data Partition, buffer *bytes.Buffer) {
	fmt.Printf("Nombre: %s, Tipo: %s, Inicio: %d, Tamaño: %d, Estado: %s, ID: %s, Ajuste: %s, Correlativo: %d\n",
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

func PrintEBR(data EBR, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Mount: %s, Fit: %s, Start: %d, Size: %d, Next: %d, Name: %s\n",
		string(data.ERBMount[:]), string(data.ERBFit[:]), data.ERBStart, data.ERBSize, data.ERBNext, string(data.ERBName[:]))
}

// =================================Estuctura MountId=================================

type MountId struct {
	MIpath   string
	MInumber int32
	MIletter int32
}
