package disco

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type mbr struct {
	Mbr_tamanio        int64
	Mbr_fecha_creacion [25]byte
	Mbr_disk_signature int64
	Disk_fit           byte
	Mbr_particiones    [4]particion
}

type particion struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
}

type ebr struct {
	Part_status byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
	Part_next   int64
}

func init() {

}

func Writefile(size int64, nombre string, path string, fit byte, unit byte) {

	var f [25]byte
	var nAleatorio int64

	_, mist := os.Stat(path)

	if os.IsNotExist(mist) {
		mist = os.Mkdir(path, 0775)
		if mist != nil {
			fmt.Println(mist)
		}
	}

	file, err := os.Create(path + "/" + nombre)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var otro int8 = 0

	s := &otro

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	escribirBytes(file, binario.Bytes())

	file.Seek(size, 0)

	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	escribirBytes(file, binario2.Bytes())

	file.Seek(0, 0)

	//Fecha actual
	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	copy(f[:], fecha)

	//Numeros aleatorios
	nAleatorio = rand.Int63n(100)

	disco := mbr{Mbr_tamanio: size, Disk_fit: fit, Mbr_fecha_creacion: f, Mbr_disk_signature: nAleatorio}

	s1 := &disco

	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())

	fmt.Println("Disco Creado")
}

func Deletefile(path string, nombre string) {
	var p string
	p = path + "/" + nombre
	p = strings.Replace(p, "\"", "", -1)
	_, mist := os.Stat(p)

	if os.IsNotExist(mist) {

	} else {
		err := os.Remove(p)
		if err != nil {
			fmt.Println("Error eliminando archivo")
		} else {
			fmt.Println("------- Archivo eliminado -------")
		}
	}
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}
