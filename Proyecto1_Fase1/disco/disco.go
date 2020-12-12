package disco

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
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

	disco.Mbr_particiones[0].Part_size = -1
	disco.Mbr_particiones[0].Part_status = 0
	disco.Mbr_particiones[1].Part_size = -1
	disco.Mbr_particiones[1].Part_status = 0
	disco.Mbr_particiones[2].Part_size = -1
	disco.Mbr_particiones[2].Part_status = 0
	disco.Mbr_particiones[3].Part_size = -1
	disco.Mbr_particiones[3].Part_status = 0

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

func CrearParticion(path string, name [16]byte, size int64, unit byte, tipo byte, fit byte) string {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return "No se encontro un disco en el path" + path
	}
	fmt.Println("aca entra")
	m := mbr{}
	var sizeMbr int64 = int64(unsafe.Sizeof(m))

	data := leerBytes(file, sizeMbr)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read dailed", err)
	}

	cPart := contarP(m)
	sumaP := cPart[0] + cPart[1]

	if sumaP == 4 {
		cadena := "----- No se pudo crear una nueva particion---\n"
		cadena += "-----La cantidad de particiones ya llegó al limite-----\n"
		cadena += "-----En este Disco: " + path
		return cadena
	} else {
		if tipo == 112 || tipo == 80 {
			fmt.Println("---- Crear Particion Primaria -----")
		} else if tipo == 80 || tipo == 101 {
			if cPart[1] == 1 {
				cadena := "----- No se puede crear otra particion extendida -----\n"
				cadena += "----- En el disco: " + path + " -----"
				return cadena
			} else if cPart[1] == 0 {
				fmt.Println("----- No hay extendidas, se puede crear una -----")
			}
		} else if tipo == 76 || tipo == 108 {
			if cPart[1] == 1 {
				fmt.Println("----- Crear logicas -----")
			} else {
				cadena := "----- ERROR - NO SE PUDO CREAR UNA NUEVA PARTICION -----\n"
				cadena += "----- No existe ninguna particion EXTENDIDA donde crear particiones LOGICAS ----- "
				cadena += "\n----- Disco: " + path + " -----"
				return cadena
			}
		}
		if !validarNombre(m, name) {
			cadena := "----- No se pudo crear la particion -----\n"
			cadena += "----- Ya existe una particion con este nombre -----\n"
			cadena += "----- En este disco: " + path + " -----"
			return cadena
		}
		disp := getEspacioDisponible(m, sizeMbr)
		sizeByte := obtenerSize(size, unit)
		if disp < int64(sizeByte) {
			cad := "-- ERROR - NO SE PUDO CREAR UNA NUEVA PARTICION --\n"
			cad += "-- No hay suficiente espacio en el disco --\n "
			cad += "-- Disponible: " + strconv.Itoa(int(disp)) + " Deseado: " + strconv.Itoa(int(sizeByte)) + " --\n"
			cad += "-- Disco: " + path
			return cad
		}

		pos := getFF(m, sizeMbr, sizeByte)

		if pos == -1 {
			cad := "-- ERROR - NO SE PUDO CREAR UNA NUEVA PARTICION --\n"
			cad += "-- No espacio suficiente ADECUADO en el disco --\n "
			cad += "-- Disco: " + path
			return cad
		}
		newParticion := particion{Part_size: int64(sizeByte)}
		newParticion.Part_status = 0
		newParticion.Part_type = tipo
		newParticion.Part_fit = fit
		newParticion.Part_name = name
		newParticion.Part_start = int64(pos)

		for i := 0; i < len(m.Mbr_particiones); i++ {
			if m.Mbr_particiones[i].Part_size == -1 {
				m.Mbr_particiones[i] = newParticion
				lastF, err := os.OpenFile(path, os.O_WRONLY, 0777)
				if err != nil {
					log.Println(err)
				}
				defer lastF.Close()
				sx := &m
				lastF.Seek(0, 0)
				var binario3 bytes.Buffer
				binary.Write(&binario3, binary.BigEndian, sx)
				escribirBytes(lastF, binario3.Bytes())
				break
			}
		}
		if newParticion.Part_type == 80 || newParticion.Part_type == 112 {
			return "--- Se creo una nueva particion PRIMARIA -- "
		} else if newParticion.Part_type == 101 || newParticion.Part_type == 69 {
			miEBR := ebr{}
			miEBR.Part_size = -1
			miEBR.Part_status = 0
			miEBR.Part_start = int64(pos)
			miEBR.Part_next = -1
			file.Seek(int64(pos), 0)
			s1 := &miEBR
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, s1)
			escribirBytes(file, binario3.Bytes())

			return "--- Se creo una nueva particion EXTENDIDA -- "
		} else if newParticion.Part_type == 76 || newParticion.Part_type == 108 {

			return "--- Se creo una nueva particion LOGICA -- "
		}
	}
	return "--- Se Llego al FINAL WHY AM I HERE ---"

}

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)

	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func contarP(m mbr) [2]int {
	numeroPrim := 0
	numeroExt := 0
	for i := 1; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_size != -1 {
			if m.Mbr_particiones[i].Part_type == 112 || m.Mbr_particiones[i].Part_type == 80 {
				numeroPrim++
			} else if m.Mbr_particiones[i].Part_type == 101 || m.Mbr_particiones[i].Part_type == 69 {
				numeroExt++
			}
		}
	}
	cant := [2]int{numeroPrim, numeroExt}
	return cant
}

func validarNombre(m mbr, nombre [16]byte) bool {
	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_name == nombre {
			return false
		}
	}
	return true
}

func getEspacioDisponible(m mbr, sizeMBR int64) int64 {
	espacioDisp := m.Mbr_tamanio - int64(sizeMBR)
	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_size > 0 {
			espacioDisp = espacioDisp - m.Mbr_particiones[i].Part_size
		}
	}
	return espacioDisp
}

func obtenerSize(num int64, unit byte) int64 {
	if unit == 107 || unit == 75 {
		return num * 1024
	} else if unit == 77 || unit == 109 {
		return num * 1024 * 1024
	} else {
		return num
	}
}

func getFF(m mbr, sizeMBR int64, sizeBuscar int64) int {
	pos := int64(sizeMBR)

	startParts := []int{0, 0, 0, 0}
	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_size > 0 {
			startParts[i] = int(m.Mbr_particiones[i].Part_start)
		}
	}
	sort.Ints(startParts)
	for i := 0; i < len(startParts); i++ {
		if startParts[i] > 0 {
			auxDisponible := startParts[i] - int(pos)
			if sizeBuscar < int64(auxDisponible) {
				return int(pos)
			}

			for j := 0; j < len(m.Mbr_particiones); j++ {
				if int(m.Mbr_particiones[j].Part_start) == startParts[i] {
					pos = int64(startParts[i]) + m.Mbr_particiones[j].Part_size
					break
				}

			}
		}
	}
	auxDisponible := int(m.Mbr_tamanio) - int(pos)
	if sizeBuscar < int64(auxDisponible) {
		return int(pos)
	}
	return -1
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func getEBR(start int64, disco *os.File) ebr {
	disco.Seek(start, 0)
	auxEB := ebr{}
	var sizeEBR int64 = int64(unsafe.Sizeof(auxEB))
	data := leerBytes(disco, sizeEBR)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &auxEB)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return auxEB
}

func CrearLogicas(path string, name [16]byte, size int64, unit byte, tipo byte, fit byte) string {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return "----- Error no se encontro el disco -----"
	}
	m := mbr{}
	var tamanioMbr int64 = int64(unsafe.Sizeof(m))

	data := leerBytes(file, tamanioMbr)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	indExt := buscarExtendida(m)
	if indExt == -1 {
		cadena := "----- Error: No se puede crear esta particion -----"
		cadena += "----- No existe ninguna particion extendida en este disco -----"
		return cadena
	}
	file.Seek(m.Mbr_particiones[indExt].Part_start, 0)
	auxEbr := ebr{}
	var sizeEBR int64 = int64(unsafe.Sizeof(auxEbr))
	data = leerBytes(file, sizeEBR)
	buffer = bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &auxEbr)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	if !validarNombrePartLogica(auxEbr, name, file) {
		cadena := "----- Error: No se pudo crear una particion -----"
		cadena += "----- Ya existe una praticion logica con este nombre en este disco -----"
		return cadena
	}
	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_name == name {
			cadena := "----- Error: No se pudo crear una particion -----"
			cadena += "----- Ya existe una praticion con este nombre en este disco -----"
			return cadena
		}
	}

	tamanioLbytes := obtenerSize(size, unit)
	disp := getEspacioEXTdisponible(m.Mbr_particiones[indExt].Part_size, sizeEBR, auxEbr, file)

	if int64(tamanioLbytes) > disp {
		cadena := "----- Error: No se pudo crear una nuevo particion logica -----"
		cadena += "----- No hay espacio suficiente en el disco -----"
		return cadena
	}
	extendidaFinal := m.Mbr_particiones[indExt].Part_start + m.Mbr_particiones[indExt].Part_size
	ultimo := getFFEXT(auxEbr, sizeEBR, tamanioLbytes, extendidaFinal, file)
	if ultimo.Part_start == -1 {
		cadena := "----- Error: No se pudo crear una nueva particion -----"
		cadena += "----- No hay espacio suficiente -----"
		return cadena
	}
	nuevoEbr := ebr{Part_status: 0}
	nuevoEbr.Part_size = int64(tamanioLbytes)
	nuevoEbr.Part_fit = fit
	nuevoEbr.Part_name = name
	nuevoEbr.Part_next = -1
	var pos int64
	var bandera bool
	if auxEbr.Part_start == ultimo.Part_start && ultimo.Part_size == -1 {
		pos = auxEbr.Part_start
	} else {
		pos = ultimo.Part_start + ultimo.Part_size
		bandera = true
	}
	nuevoEbr.Part_start = pos
	escribirEBRDisco(nuevoEbr, nuevoEbr.Part_start, file)
	if bandera {
		ultimo.Part_next = nuevoEbr.Part_start
		escribirEBRDisco(ultimo, ultimo.Part_start, file)
	}
	return "----- Se ha creado una nueva particion logica -----"
}

func buscarExtendida(m mbr) int {
	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_type == 101 {
			return i
		}
	}
	return -1
}

func validarNombrePartLogica(auxEbr ebr, name [16]byte, disco *os.File) bool {
	n := 0
	for true {
		if auxEbr.Part_name == name {
			return false
		}
		if auxEbr.Part_next == -1 {
			break
		}
		if n < 3 {
			n++
		}
		auxEbr = getEbr(auxEbr.Part_next, disco)
	}

	return true
}

func getEbr(start int64, disco *os.File) ebr {
	disco.Seek(start, 0)
	auxEbr := ebr{}
	var tamanioEbr int64 = int64(unsafe.Sizeof(auxEbr))
	data := leerBytes(disco, tamanioEbr)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &auxEbr)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	return auxEbr
}

func obtenerSizeL(tam int64, unit byte) int64 {
	if unit == 107 || unit == 75 {
		return tam * 1024
	} else if unit == 77 || unit == 109 {
		return tam * 1024 * 1024
	} else {
		return tam
	}
}

func getEspacioEXTdisponible(tamanioExt int64, tamanioEbr int64, aebr ebr, disco *os.File) int64 {
	espacioDisp := tamanioExt - int64(tamanioEbr)
	for true {
		if aebr.Part_size > 0 {
			espacioDisp = espacioDisp - aebr.Part_size
		}
		if aebr.Part_next == -1 {
			break
		}
		aebr = getEbr(aebr.Part_next, disco)
	}
	return espacioDisp
}

func getFFEXT(eb ebr, sizeEBR int64, sizeLbytes int64, finalEXT int64, disco *os.File) ebr {
	if eb.Part_next == -1 {
		return eb
	}
	for true {
		if eb.Part_next == -1 {
			break
		}
		ebrAux := getEBR(eb.Part_next, disco)
		finalActual := eb.Part_start + eb.Part_size
		disponible := ebrAux.Part_size - finalActual
		if int64(sizeLbytes) <= disponible {
			return ebrAux
		}
		eb = ebrAux
	}
	finalEBR := eb.Part_start + eb.Part_size
	disp := finalEXT - finalEBR
	if int64(sizeLbytes) <= disp {
		return eb
	}
	eb.Part_start = -1
	return eb
}

func escribirEBRDisco(nEBR ebr, posicion int64, disco *os.File) {
	disco.Seek(posicion, 0)
	s1 := &nEBR
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(disco, binario3.Bytes())
}

func BorrarParticion(path string, name [16]byte, tipoBorrado string) string {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return "----- ERROR: No se pudo abrir el disco -----"
	}

	m := mbr{}
	var tamanioMbr int64 = int64(unsafe.Sizeof(m))
	data := leerBytes(file, tamanioMbr)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	var encontrado bool = false
	var esLogica bool = false
	var ind int
	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_size != -1 {
			if m.Mbr_particiones[i].Part_name == name {
				encontrado = true
				esLogica = false
				ind = i
			}
		}
	}

	ebrAnterior := ebr{}
	auxEbr := ebr{}
	if encontrado == false {
		indExt := buscarExtendida(m)

		if indExt != -1 {
			file.Seek(m.Mbr_particiones[indExt].Part_start, 0)
			var sizeEBR int64 = int64(unsafe.Sizeof(auxEbr))
			data = leerBytes(file, sizeEBR)
			buffer = bytes.NewBuffer(data)
			err = binary.Read(buffer, binary.BigEndian, &auxEbr)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}

			ebrAnterior = auxEbr
			for true {
				if auxEbr.Part_size == -1 {
					break
				}
				if auxEbr.Part_name == name {
					encontrado = true
					esLogica = true
					break
				}

				if auxEbr.Part_next == -1 {
					break
				}
				ebrAnterior = auxEbr
				auxEbr = getEbr(auxEbr.Part_next, file)
			}
		}
	}
	if !encontrado {
		cad := "-- ERROR - NO SE PUDO BORRAR UNA PARTICION --\n"
		cad += "-- No existe ninguna particion con el nombre " + string(name[:]) + " --\n "
		cad += "-- Disco: " + path
		return cad
	}

	if esLogica {

		if tipoBorrado == "fast" {
			ebrAnterior.Part_next = auxEbr.Part_next
		} else {
			ebrAnterior.Part_next = auxEbr.Part_next

			var otro int8 = 0
			s := &otro

			file.Seek(auxEbr.Part_start, 0)
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, s)
			escribirBytes(file, binario.Bytes())

			posFinal := auxEbr.Part_start + auxEbr.Part_size
			file.Seek(posFinal, 0)
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, s)
			escribirBytes(file, binario2.Bytes())
			fmt.Println("Se Borro FULL")
		}

		file.Seek(ebrAnterior.Part_start, 0)
		s1 := &ebrAnterior

		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		escribirBytes(file, binario3.Bytes())
		return "SUCESS -- Se elimino una particion LOGICA"
	}

	if tipoBorrado == "full" {

		var otro int8 = 0
		s := &otro

		file.Seek(m.Mbr_particiones[ind].Part_start, 0)
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, s)
		escribirBytes(file, binario.Bytes())

		posFinal := m.Mbr_particiones[ind].Part_start + m.Mbr_particiones[ind].Part_size
		file.Seek(posFinal, 0)
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		escribirBytes(file, binario2.Bytes())

	}
	newParticion := particion{Part_size: -1}

	m.Mbr_particiones[ind] = newParticion

	file.Seek(0, 0)
	s1 := &m

	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())
	return "Sucess -- Se elimino una particion Primaria o extendida"

}

func AgregarAParticion(path string, name [16]byte, addValue int64, unit byte) string {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return "----- ERROR: No se pudo abrir el disco -----"
	}

	m := mbr{}
	var tamanioMbr int64 = int64(unsafe.Sizeof(m))
	data := leerBytes(file, tamanioMbr)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	var encontrado bool = false
	var tipoP int = -1
	var ind int
	for i := 0; i < len(m.Mbr_particiones); i++ {
		if m.Mbr_particiones[i].Part_size != -1 {
			if m.Mbr_particiones[i].Part_name == name {
				encontrado = true
				if m.Mbr_particiones[i].Part_type == 69 || m.Mbr_particiones[i].Part_type == 101 {
					tipoP = 1
				} else {
					tipoP = 0
				}
				ind = i
			}
		}
	}

	auxEB := ebr{}
	if encontrado == false {
		indiceEXT := buscarExtendida(m)

		if indiceEXT != -1 {
			file.Seek(m.Mbr_particiones[indiceEXT].Part_start, 0)
			var sizeEBR int64 = int64(unsafe.Sizeof(auxEB))
			data = leerBytes(file, sizeEBR)
			buffer = bytes.NewBuffer(data)
			err = binary.Read(buffer, binary.BigEndian, &auxEB)
			if err != nil {
				log.Fatal("binary.Read failed", err)
			}

			for true {
				if auxEB.Part_size == -1 {
					break
				}
				if auxEB.Part_name == name {
					encontrado = true
					tipoP = 2
					break
				}

				if auxEB.Part_next == -1 {
					break
				}
				auxEB = getEBR(auxEB.Part_next, file)
			}
		}
	}
	if !encontrado {
		cad := "----- ERROR: No se pudo agregar espacio a la particion -----\n"
		cad += "----- No existe ninguna particion con este nombre -----"
		return cad
	}

	addValue = obtenerSize(addValue, unit)

	if tipoP == 0 {

		nuevoSize := m.Mbr_particiones[ind].Part_size + int64(addValue)
		nuevoFinal := m.Mbr_particiones[ind].Part_start + nuevoSize
		if nuevoFinal > m.Mbr_tamanio || nuevoSize <= 0 {
			cad := "----- ERROR: No se pudo agregar espacio a la particion -----\n"
			return cad
		}

		for i := 0; i < len(m.Mbr_particiones); i++ {
			if (m.Mbr_particiones[i].Part_size != -1) && m.Mbr_particiones[i].Part_name != m.Mbr_particiones[ind].Part_name {
				auxFinal := m.Mbr_particiones[i].Part_start + m.Mbr_particiones[i].Part_size
				auxInicio := m.Mbr_particiones[i].Part_start

				if nuevoFinal > auxInicio && nuevoFinal < auxFinal {
					cad := "----- Error: No se pudo añadir a una particion -----\n"
					cad += "----- El tamaño es muy grande -----"
					return cad
				}

			}
		}

		fmt.Print("Tamaño anterior: " + string(m.Mbr_particiones[ind].Part_size))
		fmt.Print("Tamaño para agregar: " + string(nuevoSize))
		m.Mbr_particiones[ind].Part_size = nuevoSize
		fmt.Print("Tamaño total: " + string(m.Mbr_particiones[ind].Part_size))

		file.Seek(0, 0)
		s1 := &m

		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		escribirBytes(file, binario3.Bytes())
		return "----- Se agrego espacio a una particion primaria -----"

	} else if tipoP == 1 {

		nuevoSize := m.Mbr_particiones[ind].Part_size + int64(addValue)
		nuevoFinal := m.Mbr_particiones[ind].Part_start + nuevoSize
		if nuevoFinal > m.Mbr_tamanio || nuevoSize <= 0 {
			cad := "----- ERROR : No se pudo agregar mas espacio a la particion -----"
			return cad
		}

		for i := 0; i < len(m.Mbr_particiones); i++ {

			if (m.Mbr_particiones[i].Part_size != -1) && m.Mbr_particiones[i].Part_name != m.Mbr_particiones[ind].Part_name {
				auxFinal := m.Mbr_particiones[i].Part_start + m.Mbr_particiones[i].Part_size
				auxInicio := m.Mbr_particiones[i].Part_start

				if nuevoFinal > auxInicio && nuevoFinal < auxFinal {
					cad := "----- ERROR : No se pudo agregar mas espacio a la particion -----\n"
					cad += "----- El tamaño es muy grande ----- "
					return cad
				}

			}
		}

		e := ebr{}
		file.Seek(m.Mbr_particiones[ind].Part_start, 0)
		var sizeEBR int64 = int64(unsafe.Sizeof(e))
		data = leerBytes(file, sizeEBR)
		buffer = bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &e)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		for true {

			if e.Part_size == -1 {
				break
			}
			auxF := e.Part_size + e.Part_start
			if nuevoFinal > e.Part_start && nuevoFinal < auxF {
				cad := "----- ERROR : No se pudo agregar mas espacio a la particion -----\n"
				cad += "----- El tamaño es muy grande -----  "
				return cad
			}
			if e.Part_next == -1 {
				break
			}
			e = getEBR(e.Part_next, file)
		}

		fmt.Print("Tamaño anterior: " + string(m.Mbr_particiones[ind].Part_size))
		fmt.Print("Tamaño para agregar: " + string(nuevoSize))
		m.Mbr_particiones[ind].Part_size = nuevoSize
		fmt.Print("Tamaño total: " + string(m.Mbr_particiones[ind].Part_size))

		file.Seek(0, 0)
		s1 := &m

		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		escribirBytes(file, binario3.Bytes())
		return "----- Se agrego espacio a una particion  -----"

	} else if tipoP == 2 {

		nuevoSize := auxEB.Part_size + int64(addValue)
		indiceEXT := buscarExtendida(m)
		ex := ebr{}
		file.Seek(m.Mbr_particiones[indiceEXT].Part_start, 0)
		var sizeEBR int64 = int64(unsafe.Sizeof(ex))
		data = leerBytes(file, sizeEBR)
		buffer = bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &ex)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		var disponibleEXT int64 = DispEXT(m.Mbr_particiones[indiceEXT].Part_size, ex, file)
		nuevoFinal := auxEB.Part_start + nuevoSize
		if nuevoFinal > m.Mbr_tamanio || nuevoSize <= 0 {
			cad := "----- ERROR : No se pudo agregar mas espacio a la particion -----\n"
			cad += "----- El tamaño es muy grande o es negativo -----"
			return cad
		}

		disponibleEXT += auxEB.Part_size
		if nuevoSize > disponibleEXT {
			cad := "----- ERROR : No se pudo agregar mas espacio a la particion -----\n"
			cad += "----- El tamaño es muy grande -----"
			return cad
		}

		for true {

			if ex.Part_size == -1 {
				break
			}
			if ex.Part_name != auxEB.Part_name {
				auxF := ex.Part_size + ex.Part_start
				if nuevoFinal > ex.Part_start && nuevoFinal < auxF {
					cad := "----- ERROR : No se pudo agregar mas espacio a la particion -----\n"
					cad += "----- El tamaño es muy grande ----- "
					return cad
				}
			}
			if ex.Part_next == -1 {
				break
			}
			ex = getEBR(ex.Part_next, file)
		}

		fmt.Print("Tamaño anterior: " + string(m.Mbr_particiones[ind].Part_size))
		fmt.Print("Tamaño para agregar: " + string(nuevoSize))
		m.Mbr_particiones[ind].Part_size = nuevoSize
		fmt.Print("Tamaño total: " + string(m.Mbr_particiones[ind].Part_size))

		file.Seek(auxEB.Part_start, 0)
		s1 := &auxEB

		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, s1)
		escribirBytes(file, binario3.Bytes())
		return "----- Se agrego espacio a una particion logica -----"
	} else {
		return "ERROR"
	}
}

func DispEXT(sizeTotal int64, ext ebr, disco *os.File) int64 {
	disp := sizeTotal
	for true {
		if ext.Part_size == -1 {
			break
		}
		disp -= ext.Part_size
		if ext.Part_next == -1 {
			break
		}
		ext = getEBR(ext.Part_next, disco)
	}
	return disp
}

func MBR(destino string, pathDisco string) bool {
	disco, err := os.OpenFile(pathDisco, os.O_RDWR, 0777)
	defer disco.Close()
	if err != nil {
		fmt.Println("ERROR - No se pudo abrir el disco para obtener el MBR")
		return false
	}

	m := mbr{}
	var sizeMBR int64 = int64(unsafe.Sizeof(m))
	//Leer la cantidad 'sizeMBR' de bytes sobre el archivo abierto.
	data := leerBytes(disco, sizeMBR)
	buffer := bytes.NewBuffer(data)
	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	escribirGraficaMBR(m, destino, disco)
	compilarGraficaMBR(destino, "./Dots/MBR.dot")
	return true
}

func compilarGraficaMBR(destino string, from string) {
	carpetas := strings.Split(destino, "/")
	ruta := "/home"
	for i := 2; i < len(carpetas)-1; i++ {
		ruta += "/" + carpetas[i]
	}
	os.MkdirAll(ruta, 0777)
	p, _ := exec.LookPath("dot")
	formato := strings.Split(destino, ".")
	cmd, _ := exec.Command(p, "-Tpng", from).Output()
	if formato[1] == "jpg" || formato[1] == "jpeg" {
		cmd, _ = exec.Command(p, "-Tjpg", from).Output()
	} else if formato[1] == "png" {
		cmd, _ = exec.Command(p, "-Tpng", from).Output()
	} else {
		cmd, _ = exec.Command(p, "-Tpdf", from).Output()
	}

	mode := int(0777)
	ioutil.WriteFile(destino, cmd, os.FileMode(mode))
}

func escribirGraficaMBR(m mbr, destino string, disco *os.File) {
	file, err := os.Create("./Dots/MBR.dot")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var cadena string = ""
	cadena += "digraph MBR{\n"
	cadena += "    MBR[\n"
	cadena += "    shape=plaintext\n"
	cadena += "    label=<\n"
	cadena += "    <table border='1' cellborder='1'>\n"
	cadena += "    <tr><td>Nombre</td><td>Valor</td></tr>\n"
	//Escribir informacion del MBR+
	cadena += "	   <tr><td>MBR_Tamanio</td><td>" + strconv.Itoa(int(m.Mbr_tamanio)) + "</td></tr>\n"
	cadena += "    <tr><td>MBR_Fecha_Creacion</td><td>" + dateToGoString(m.Mbr_fecha_creacion) + "</td></tr>\n"
	cadena += "    <tr><td>MBR_Disk_Signature</td><td>" + strconv.Itoa(int(m.Mbr_disk_signature)) + "</td></tr>\n"
	//Recorrer Particiones
	iExt := -1
	for i := 0; i < len(m.Mbr_particiones); i++ {
		n := strconv.Itoa(i + 1)
		tipo := m.Mbr_particiones[i].Part_type
		if m.Mbr_particiones[i].Part_size == -1 {
			cadena += "<tr><td>part" + n + "_Estado</td><td>-1</td></tr>\n"
			cadena += "<tr><td>part" + n + "_Type</td><td> - </td></tr>\n"
			cadena += "<tr><td>part" + n + "_Fit</td><td> -</td></tr>\n"
			cadena += "<tr><td>part" + n + "_Start</td><td> - </td></tr>\n"
			cadena += "<tr><td>part" + n + "_Name</td><td> - </td></tr>\n"
			cadena += "<tr><td>part" + n + "_Size</td><td> - </td></tr>\n"
		} else {
			stat := m.Mbr_particiones[i].Part_status
			mystat := string(stat)
			if stat == 2 {
				mystat = "2"
			} else if stat == 1 {
				mystat = "1"
			} else {
				mystat = "0"
			}
			cadena += "<tr><td>part" + n + "_Estado</td><td>" + mystat + "</td></tr>\n"
			cadena += "<tr><td>part" + n + "_Type</td><td>" + string(tipo) + "</td></tr>\n"
			cadena += "<tr><td>part" + n + "_Fit</td><td>" + string(m.Mbr_particiones[i].Part_fit) + "</td></tr>\n"
			cadena += "<tr><td>part" + n + "_Start</td><td>" + strconv.Itoa(int(m.Mbr_particiones[i].Part_start)) + "</td></tr>\n"
			cadena += "<tr><td>part" + n + "_Name</td><td>" + nameToGoString(m.Mbr_particiones[i].Part_name) + "</td></tr>\n"
			cadena += "<tr><td>part" + n + "_Size</td><td>" + strconv.Itoa(int(m.Mbr_particiones[i].Part_size)) + "</td></tr>\n"
		}
		if tipo == 101 || tipo == 69 { //Si es extendida guardar el indice
			iExt = i
		}
	}
	if iExt != -1 {
		//Hay extendida entonces agregar los EBR a la grafica
		st := m.Mbr_particiones[iExt].Part_start
		eb := getEBR(st, disco)
		x := 1
		for true {
			if eb.Part_size == 1 {
				break
			}
			num := strconv.Itoa(x)
			stat := eb.Part_status
			mystat := string(stat)
			if stat == 2 {
				mystat = "2"
			} else if stat == 1 {
				mystat = "1"
			} else {
				mystat = "0"
			}
			cadena += "<tr><td>EBR" + num + "_Estado</td><td>" + mystat + "</td></tr>\n"
			cadena += "<tr><td>EBR" + num + "_Fit</td><td>" + string(eb.Part_fit) + "</td></tr>\n"
			cadena += "<tr><td>EBR" + num + "_Start</td><td>" + strconv.Itoa(int(eb.Part_start)) + "</td></tr>\n"
			cadena += "<tr><td>EBR" + num + "_Size</td><td>" + strconv.Itoa(int(eb.Part_size)) + "</td></tr>\n"
			cadena += "<tr><td>EBR" + num + "_Name</td><td>" + nameToGoString(eb.Part_name) + "</td></tr>\n"
			cadena += "<tr><td>EBR" + num + "_Next</td><td>" + strconv.Itoa(int(eb.Part_next)) + "</td></tr>\n"
			if eb.Part_next == -1 {
				break
			}
			x++
			eb = getEBR(eb.Part_next, disco)
		}
	}

	cadena += "    </table>\n"
	cadena += "    >];\n"
	cadena += "}"
	file.Sync() // flush
	buffer := bufio.NewWriter(file)
	buffer.WriteString(cadena)
	if err != nil {
		panic(err)
	}

	buffer.Flush() // Flush writes any buffered .
}

func dateToGoString(c [25]byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func nameToGoString(c [16]byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}
