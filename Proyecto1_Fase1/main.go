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
	fmt.Println("Saliendo De La Terminal......")
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
				fmt.Println(s.Text())
				lineaComando(strings.TrimSpace(s.Text()))
			}

		}
	}
}

//Mkdisk -Size->5 -unit->K -path->/home/user/Disco1.dsk
//rmDisk -path->/home/user/Disco1.dsk
//mkdisk -size->5 -unit->K -path->"/home/mis discos/Disco3.dsk"
//rmdisk -path->"/home/mis discos/Disco3.dsk"
