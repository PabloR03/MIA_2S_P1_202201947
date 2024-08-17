package Estructura

import (
	"bytes"
	"fmt"
)

//  =============================================================

type MRB struct {
	MbrSize      int32
	CreationDate [10]byte
	Signature    int32
	Fit          [1]byte
	Partitions   [4]Partition
}

func PrintMBR(data MRB, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Size: %d, CreationDate: %s, Signature: %d, Fit: %s\n", data.MbrSize, string(data.CreationDate[:]), data.Signature, string(data.Fit[:]))
	for i := 0; i < 4; i++ {
		PrintPartition(data.Partitions[i], buffer)
	}
}

//  =============================================================

type Partition struct {
	Start       int32
	Correlative int32
	Size        int32
	Unit        [1]byte
	Path        [100]byte
	Type        [1]byte
	Fit         [2]byte
	Name        [16]byte
	Status      [1]byte // Puede ser '1' para activa y '0' para inactiva, según se necesite
	PartId      [4]byte
}

func PrintPartition(data Partition, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Start: %d, Correlative: %d, Size: %d, Unit: %s, Path: %s, Type: %s, Fit: %s, Name: %s, Status: %s\n",
		data.Start, data.Correlative, data.Size, string(data.Unit[:]), string(data.Path[:]), string(data.Type[:]),
		string(data.Fit[:]), string(data.Name[:]), string(data.Status[:]))
}

//  =============================================================

type EBR struct {
	Part_Mount [1]byte
	Part_Fit   [1]byte
	Part_Start int32
	Part_Size  int32
	Part_Next  int32
	Part_Name  [16]byte
}

func PrintEBR(data EBR, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "Mount: %s, Fit: %s, Start: %d, Size: %d, Next: %d, Name: %s\n",
		string(data.Part_Mount[:]), string(data.Part_Fit[:]), data.Part_Start, data.Part_Size, data.Part_Next, string(data.Part_Name[:]))
}

// =============================================================

type MountId struct {
	MIRuta   string
	MINumero int32
	MILetra  int32
}
