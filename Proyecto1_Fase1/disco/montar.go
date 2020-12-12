package disco

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func init() {
}

type Mount struct {
	Letra     byte
	Path      [100]byte
	Status    uint8
	Particion [50]ParticionMontada
}

type ParticionMontada struct {
	Id     [10]byte
	Name   [16]byte
	Tipo   byte
	Start  int64
	Status byte
}

func MountDisk(path string, name [16]byte, letra byte) Mount {
	newMount := Mount{}
	newMount.Letra = letra
	var binPath [100]byte
	copy(binPath[:], path)
	newMount.Path = binPath
	newMount.Status = 1
	newMount = llenarParts(newMount)
	newMount.Particion[0] = MountParticion(path, name, letra, 1)
	if newMount.Particion[0].Status == 101 {
		fmt.Println("ERROR 101 - No se monto la particion")
		return Mount{Status: 101}
	}
	return newMount
}

func llenarParts(newMount Mount) Mount {
	for i := 0; i < len(newMount.Particion); i++ {
		newMount.Particion[i].Status = 90
	}
	return newMount
}
func MountParticion(path string, name [16]byte, letra byte, numero int) ParticionMontada {
	nuevaParticion := ParticionMontada{}
	nuevaParticion.Name = name
	aux := string(letra)
	miId := "vd" + aux + strconv.Itoa(numero)
	var bId [10]byte
	copy(bId[:], miId)
	nuevaParticion.Id = bId
	nuevaParticion.Status = 0

	disco, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer disco.Close()
	if err != nil {
		fmt.Println("Error, no existe este disco en esta ruta")
		return ParticionMontada{Status: 101}
	}
	m := mbr{}
	var tamanioMbr int64 = int64(unsafe.Sizeof(m))
	data := leerBytes(disco, tamanioMbr)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	var bandera bool = false

	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_name == name {
			m.Mbr_particiones[i].Part_status = 1
			nuevaParticion.Tipo = 0
			nuevaParticion.Start = m.Mbr_particiones[i].Part_start
			bandera = true
			disco.Seek(0, 0)
			s1 := &m

			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s1)
			escribirBytes(disco, binario3.Bytes())
			break
		}
	}
	if bandera == false {
		var iExt int = -1
		for i := 0; i < len(m.Mbr_particiones); i++ {
			if m.Mbr_particiones[i].Part_type == 101 || m.Mbr_particiones[i].Part_type == 69 {
				iExt = i
				break
			}
		}
		e := ebr{}
		disco.Seek(m.Mbr_particiones[iExt].Part_start, 0)
		var tamanioEbr int64 = int64(unsafe.Sizeof(e))
		data = leerBytes(disco, tamanioEbr)
		buffer = bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &e)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		for true {
			if e.Part_size == -1 {
				break
			}
			if e.Part_name == name {
				e.Part_status = 1
				nuevaParticion.Tipo = 1
				nuevaParticion.Start = e.Part_start
				bandera = true
				disco.Seek(e.Part_start, 0)
				s1 := &e
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, s1)
				escribirBytes(disco, binario3.Bytes())
				break
			}
			if e.Part_next == -1 {
				break
			}
			e = getEBR(e.Part_next, disco)
		}
	}
	if bandera == false {
		fmt.Println("----- Error: No se encontro una particion con ese nombre en esa particion -----")
		return ParticionMontada{Status: 101}
	}
	for c := 0; c < len(m.Mbr_particiones); c++ {
		fmt.Println(m.Mbr_particiones[c].Part_status)
	}
	return nuevaParticion
}

func AgregarParticion(p ParticionMontada, disco Mount) Mount {
	if p.Status == 101 {
		return Mount{Status: 101}
	}
	for i := 0; i < len(disco.Particion); i++ {
		if disco.Particion[i].Status == 90 {
			//Colocar en este indice
			disco.Particion[i] = p
			return disco
		}
	}
	fmt.Println("ERROR - El disco ya esta lleno de particions montadas")
	return Mount{Status: 101}
}

func EstaMontada(disc Mount, name [16]byte) bool {
	for i := 0; i < len(disc.Particion); i++ {
		if disc.Particion[i].Name == name && disc.Particion[i].Status != 90 {
			return true
		}
	}
	return false
}

func GenerarLetra(indice int) byte {
	var ascii int = 97 + indice
	letra := byte(ascii)
	return letra
}

func GenerarNum(disc Mount) int {
	for i := 0; i < len(disc.Particion); i++ {
		if disc.Particion[i].Status == 90 {
			return i + 1
		}
	}
	return -1
}

func Buscar(id string) string {
	var binid [10]byte
	var mounts [26]Mount
	copy(binid[:], id)
	for j := 0; j < len(mounts); j++ {
		fmt.Println(mounts[j].Status)
		if mounts[j].Status != 0 {
			for k := 0; k < len(mounts[j].Particion); k++ {
				if mounts[j].Particion[k].Status != 90 {
					if mounts[j].Particion[k].Id == binid {
						mounts[j].Particion[k].Status = 90
						fmt.Println("----- Se Desmonto -----")
						s := string(mounts[j].Path[:])
						return s
					}
				}
			}
		}
	}
	return "ERROR"
}
