package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	d "./disco"
)

var mounts [26]d.Mount

func main() {
	verificacion()
}

func verificacion() {
	fmt.Println("----- Terminal -----")
	salir := 0

	for salir != 1 {
		fmt.Println("Ingrese Comando: ")
		reader := bufio.NewReader(os.Stdin)
		command, _, _ := reader.ReadLine()
		if string(command) == "x" {
			salir = 1
		} else if strings.Contains(string(command), "#") {
			fmt.Println(string(command))
		} else {
			if string(command) != "" {
				lineaComando(string(command))
			}
		}
	}
	fmt.Println("----- Finalizando ejecucion -----")
}

func lineaComando(comando string) {
	var commandArray []string
	commandArray = strings.Split(comando, " ")
	ejecutarComando(commandArray)
}

func ejecutarComando(commandArray []string) {
	data := strings.ToLower(commandArray[0])
	if strings.ToLower(data) == "mkdisk" {
		mkdisk(commandArray)
	} else if strings.ToLower(data) == "rmdisk" {
		rmdisk(commandArray)
	} else if strings.ToLower(data) == "pause" {
		pause()
	} else if strings.ToLower(data) == "exec" {
		exec(commandArray)
	} else if strings.ToLower(data) == "fdisk" {
		fdisk(commandArray)
	} else if strings.ToLower(data) == "rep" {
		rep(commandArray)
	} else if strings.ToLower(data) == "mount" {
		mount(commandArray)
	} else {
		fmt.Println("Otro Comando")
	}
}

func mkdisk(command []string) {
	var commandArray []string
	var path, nombre string
	var size int64
	var fit, unit byte
	var vPath, vSize, vFit, vUnit, errFit, errUnit bool

	fit = 'F'
	unit = 'M'

	for i := 1; i < len(command); i++ {
		commandArray = strings.Split(command[i], ">")
		for a := 0; a < len(commandArray); a++ {
			commandArray[a] = strings.Trim(commandArray[a], "-")
		}

		if strings.ToLower(commandArray[0]) == "size" {
			n, err := strconv.ParseInt(commandArray[1], 10, 64)
			if err == nil {
				size = n * 1024 * 1024
			}
			vSize = true
		} else if strings.ToLower(commandArray[0]) == "path" {
			if strings.Contains(commandArray[1], "\"") {
				path = commandArray[1] + " " + command[i+1]
				i++
				fmt.Println(path)
				for s := i; s < len(command); s++ {
					if strings.Contains(command[s], "\"") {
						path = path + " " + command[s]
						i++
						break
					} else {
						path = path + " " + command[s]
						i++
					}
				}
				path = strings.Trim(path, "\"")
				vPath = true
			} else {
				path = commandArray[1]
				vPath = true
			}
			arrayPath := strings.Split(path, "/")
			path = ""
			for q := 1; q < len(arrayPath); q++ {
				if strings.Contains(arrayPath[q], ".dsk") {
					nombre = arrayPath[q]
				} else {
					path = path + "/" + arrayPath[q]
				}

			}
		} else if strings.ToLower(commandArray[0]) == "unit" {
			if strings.ToUpper(commandArray[1]) == "M" {
				unit = 'M'
				errUnit = false
			} else if strings.ToUpper(commandArray[1]) == "K" {
				unit = 'K'
				size = size / 1024
				errUnit = false
			} else {
				unit = 'N'
				errUnit = true
				fmt.Println("Unidad no valida")
			}
			vUnit = true
		} else if strings.ToLower(commandArray[0]) == "fit" {
			if strings.ToUpper(commandArray[1]) == "BF" {
				f := byte('B')
				fit = f
				errFit = false
			} else if strings.ToUpper(commandArray[1]) == "FF" {
				f := byte('F')
				fit = f
				errFit = false
			} else if strings.ToUpper(commandArray[1]) == "WF" {
				f := byte('W')
				fit = f
				errFit = false
			} else {
				fmt.Println("Ajuste incorrecto")
				f := byte('E')
				fit = f
				errFit = true
			}
			vFit = true
		}
	}

	if vPath == true && vSize == true {
		if vFit == true || vUnit == true {
			if errFit == true {
				fmt.Println("Error en el parametro Fit")
			} else if errUnit == true {
				fmt.Println("Errore en el parametro Unit")
			} else {
				d.Writefile(size, nombre, path, fit, unit)
			}
		}
	} else {
		if vPath == false && vSize == false {
			fmt.Println("Falta Path y Size")
		} else if vPath == false && vSize == true {
			fmt.Println("Falta Path")
		} else if vPath == true && vSize == false {
			fmt.Println("Falta Size")
		}
	}
}

func rmdisk(command []string) {
	var commandArray []string
	var path, nombre string
	var vPath = false
	for i := 1; i < len(command); i++ {
		commandArray = strings.Split(command[i], ">")
		for a := 0; a < len(commandArray); a++ {
			commandArray[a] = strings.Trim(commandArray[a], "-")
		}
		if strings.ToLower(commandArray[0]) == "path" {
			if strings.Contains(commandArray[1], "\"") {
				path = commandArray[1] + " " + command[i+1]
				i++
				for s := i; s < len(command); s++ {
					if strings.Contains(command[s], "\"") {
						path = path + " " + command[s]
						i++
						break
					} else {
						path = path + " " + command[s]
						i++
					}
					vPath = true
				}
			} else {
				path = commandArray[1]
				vPath = true
			}
			arrayPath := strings.Split(path, "/")
			path = ""
			for q := 1; q < len(arrayPath); q++ {
				if strings.Contains(arrayPath[q], ".dsk") {
					nombre = arrayPath[q]
				} else {
					path = path + "/" + arrayPath[q]
				}

			}
		}
	}
	fmt.Println(vPath)
	d.Deletefile(path, nombre)
}

func pause() {
	salir := 0
	for salir != 1 {
		fmt.Println("------Presione X para continuar ------")
		reader := bufio.NewReader(os.Stdin)
		command, _, _ := reader.ReadLine()
		if strings.ToLower(string(command)) == "x" {
			salir = 1
		}
	}
}

func exec(command []string) {
	var path string = ""
	var commandArray []string
	for i := 1; i < len(command); i++ {
		commandArray = strings.Split(command[i], ">")
		for a := 0; a < len(commandArray); a++ {
			commandArray[a] = strings.Trim(commandArray[a], "-")
		}

		if strings.ToLower(commandArray[0]) == "path" {
			path = commandArray[1]

			fmt.Println(path)
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			s := bufio.NewScanner(file)
			for s.Scan() {
				//fmt.Println(s.Text())
				if strings.Contains(string(s.Text()), "#") {
					fmt.Println(string(s.Text()))
				} else {
					if string(s.Text()) != "" {
						lineaComando(string(s.Text()))
					}
				}
			}
		}
	}
}

func fdisk(command []string) {
	var commandArray []string
	var size, añadir int64
	var unit, tipo, fit byte
	var path, nombre, delete string
	var S, P, A, D, N, errUnit, errFit bool
	var name [16]byte

	unit = 'M'
	tipo = 'P'
	fit = 'W'

	for i := 1; i < len(command); i++ {
		commandArray = strings.Split(command[i], ">")
		for a := 0; a < len(commandArray); a++ {
			commandArray[a] = strings.Trim(commandArray[a], "-")
		}
		if strings.ToLower(commandArray[0]) == "size" && D == false && A == false {
			n, err := strconv.ParseInt(commandArray[1], 10, 64)
			if err == nil {
				size = n
			}
			S = true
		} else if strings.ToLower(commandArray[0]) == "unit" {
			if strings.ToUpper(commandArray[1]) == "M" {
				unit = 'M'
				errUnit = false
			} else if strings.ToUpper(commandArray[1]) == "K" {
				unit = 'K'
				errUnit = false
			} else {
				unit = 'N'
				errUnit = true
				fmt.Println("Unidad no valida")
			}
		} else if strings.ToLower(commandArray[0]) == "path" {
			if strings.Contains(commandArray[1], "\"") {
				path = commandArray[1] + " " + command[i+1]
				i++
				fmt.Println(path)
				for s := i; s < len(command); s++ {
					if strings.Contains(command[s], "\"") {
						path = path + " " + command[s]
						i++
						break
					} else {
						path = path + " " + command[s]
						i++
					}
				}
				path = strings.Trim(path, "\"")
			} else {
				path = commandArray[1]
			}
			arrayPath := strings.Split(path, "/")
			path = ""
			for q := 1; q < len(arrayPath); q++ {
				if strings.Contains(arrayPath[q], ".dsk") {
					nombre = arrayPath[q]
				} else {
					path = path + "/" + arrayPath[q]
				}
			}
			P = true
		} else if strings.ToLower(commandArray[0]) == "type" {
			if strings.ToUpper(commandArray[1]) == "P" {
				tipo = 'P'
			} else if strings.ToUpper(commandArray[1]) == "E" {
				tipo = 'E'
			} else if strings.ToUpper(commandArray[1]) == "L" {
				tipo = 'L'
			} else {
				tipo = 'N'
			}
		} else if strings.ToLower(commandArray[0]) == "fit" {
			if strings.ToUpper(commandArray[1]) == "BF" {
				f := byte('B')
				fit = f
				errFit = false
			} else if strings.ToUpper(commandArray[1]) == "FF" {
				f := byte('F')
				fit = f
				errFit = false
			} else if strings.ToUpper(commandArray[1]) == "WF" {
				f := byte('W')
				fit = f
				errFit = false
			} else {
				fmt.Println("Ajuste incorrecto")
				f := byte('E')
				fit = f
				errFit = true
			}
		} else if strings.ToLower(commandArray[0]) == "delete" && S == false && A == false {
			D = true
			if strings.ToLower(commandArray[1]) == "fast" {
				delete = "fast"
			} else if strings.ToLower(commandArray[1]) == "full" {
				delete = "full"
			}
		} else if strings.ToLower(commandArray[0]) == "name" {
			copy(name[:], commandArray[1])
			N = true
		} else if strings.ToLower(commandArray[0]) == "add" && S == false && A == false {
			n, err := strconv.ParseInt(commandArray[1], 10, 64)
			if err == nil {
				añadir = n
			}
			A = true
		}

	}
	path = path + "/" + nombre

	if S == true {
		if !errFit || !errUnit {
			if P == true && N == true {
				if tipo == 'L' {
					res := d.CrearLogicas(path, name, size, unit, tipo, fit)
					fmt.Println(res)
				} else {
					res := d.CrearParticion(path, name, size, unit, tipo, fit)
					fmt.Println(res)
				}
			} else {
				fmt.Println("Falta algun parametro obligatorio")
			}
		}
	} else if D == true {
		if !errFit || !errUnit {
			var bName [16]byte
			bName = name
			if strings.ToLower(delete) != "fast" && strings.ToLower(delete) != "full" {
				fmt.Println("----- Error: Tipo invalido -----")
			} else {
				res := d.BorrarParticion(path, bName, strings.ToLower(delete))
				fmt.Println(res)
			}
		}
	} else if A == true {
		if !errFit || !errUnit {
			res := (d.AgregarAParticion(path, name, añadir, unit))
			fmt.Println(res)
		}
	}
}

func rep(command []string) {
	var name, path, id string
	var n, p, ident bool
	var commandArray []string
	for i := 1; i < len(command); i++ {
		commandArray = strings.Split(command[i], ">")
		for a := 0; a < len(commandArray); a++ {
			commandArray[a] = strings.Trim(commandArray[a], "-")
		}
		if strings.ToLower(commandArray[0]) == "name" {
			name = commandArray[1]
			n = true
		} else if strings.ToLower(commandArray[0]) == "path" {
			if strings.Contains(commandArray[1], "\"") {
				path = commandArray[1] + " " + command[i+1]
				i++
				for s := i; s < len(command); s++ {
					if strings.Contains(command[s], "\"") {
						path = path + " " + command[s]
						i++
						break
					} else {
						path = path + " " + command[s]
						i++
					}
				}
			} else {
				path = commandArray[1]
			}
			p = true
		} else if strings.ToLower(commandArray[0]) == "id" {
			id = commandArray[1]
			ident = true
		}
	}
	fmt.Println(path)
	fmt.Println(name)
	fmt.Println(id)
	fmt.Println(p)
	fmt.Println(n)
	fmt.Println(ident)
	if n == true && p == true && ident == true {
		if strings.ToLower(name) == "mbr" {
			var binid [10]byte
			var s string
			copy(binid[:], id)
			for j := 0; j < len(mounts); j++ {
				//fmt.Println(mounts[j].Status)
				if mounts[j].Status != 0 {
					for k := 0; k < len(mounts[j].Particion); k++ {
						if mounts[j].Particion[k].Status != 90 {
							if mounts[j].Particion[k].Id == binid {
								mounts[j].Particion[k].Status = 90
								fmt.Println("----- Se Desmonto -----")
								s = string(mounts[j].Path[:])
							}
						}
					}
				}
			}
			if s != "ERROR" {
				t := d.MBR(strings.TrimSpace(path), strings.TrimSpace(s))
				if t {
					fmt.Println("----- Se grafico exitosamente -----")
				} else {
					fmt.Println("----- Ocurrio un error -----")
				}
			} else {
				fmt.Println("----- Particion no montada -----")
			}
		} else if strings.ToLower(name) == "disk" {

		} else {
			fmt.Println("----- Error: Tipo de reporte equivocado -----")
		}
	} else {
		fmt.Println("----- Error: hace falta un parametro -----")
	}

}

func mount(command []string) {
	var commandArray []string
	var path, name string
	var p, n bool

	for i := 1; i < len(command); i++ {
		commandArray = strings.Split(command[i], ">")
		for a := 0; a < len(commandArray); a++ {
			commandArray[a] = strings.Trim(commandArray[a], "-")
		}
		if strings.ToLower(commandArray[0]) == "name" {
			name = commandArray[1]
			n = true
		} else if strings.ToLower(commandArray[0]) == "path" {
			if strings.Contains(commandArray[1], "\"") {
				path = commandArray[1] + " " + command[i+1]
				i++
				for s := i; s < len(command); s++ {
					if strings.Contains(command[s], "\"") {
						path = path + " " + command[s]
						i++
						break
					} else {
						path = path + " " + command[s]
						i++
					}
				}
			} else {
				path = commandArray[1]
			}
			p = true
		}
	}
	if n == true && p == true {
		fmt.Println("Montar Particion: " + name + " ubicada en: " + path)
		var bName [16]byte
		copy(bName[:], name)
		var bPath [100]byte
		copy(bPath[:], path)
		var h bool = false
		for i := 0; i < len(mounts); i++ {
			if mounts[i].Status == 1 {
				if d.EstaMontada(mounts[i], bName) {
					fmt.Println("----- Error, esta particion ya está montada -----")
					h = true
					break
				} else {
					n := d.GenerarNum(mounts[i])
					fmt.Println("N: ", n)
					NuevaParticion := d.MountParticion(path, bName, mounts[i].Letra, n)
					NuevoMontado := d.AgregarParticion(NuevaParticion, mounts[i])
					if NuevoMontado.Status == 101 {
						fmt.Println("----- No se pudo montar ningun disco -----")
						break
					}
					mounts[i] = NuevoMontado
					fmt.Println("----- Se logro montar la particion -----")
					h = true
				}
			}
		}
		if h == false {
			//Este disco no esta montado, montarlo y asignarle letra
			for i := 0; i < len(mounts); i++ {
				if mounts[i].Status == 0 {
					//Montar en este indice
					letra := d.GenerarLetra(i)
					NuevoMontado := d.MountDisk(path, bName, letra)
					if NuevoMontado.Status == 101 {
						break
					}
					mounts[i] = NuevoMontado
					break
				}
			}
		}
	} else {
		fmt.Println("----- Error: hace falta un parametro -----")
	}
}

//Mkdisk -Size->5 -unit->M -path->/home/Disco1.dsk
//rmDisk -path->/home/Disco1.dsk
//mkdisk -size->5 -unit->K -path->"/home/mis discos/Disco3.dsk"
//rmdisk -path->"/home/mis discos/Disco3.dsk"
//fdisk -Size->1 -add->5 -path->/home/Disco1.dsk -name->Particion1
//mount -path->/home/Disco1.dsk -name->Particion1
//rep -id->vda1 -Path->/home/user/reports/reporte1.jpg -name->mbr
