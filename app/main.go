package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// creamos una estructura para almacenar los datos
type Email struct {
	messageId string
	date      string
	from      string
	to        string
	subject   string
	body      string
}

func main() {
	//iniciar servidor de profiling en un goroutine
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	rootDir, err := os.Open("./allen-p")
	if err != nil {
		log.Fatal(err)
	}

	// Llamar a la función recursiva
	findsDir(rootDir)

	// Cerrar el directorio raíz
	// Forzar a vaciar el búfer de salida
	os.Stdout.Sync()

	// Bucle infinito para mantener el programa en ejecución
	select {}
}

func findsDir(dir *os.File) {
	// Leer el directorio raíz
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	// Recorrer los ficheros del directorio raíz
	for _, fi := range fileInfos {
		// Comprobar si es un directorio
		if fi.IsDir() {
			// Si es un directorio, abrirlo
			subDir, err := os.Open(dir.Name() + "/" + fi.Name())
			if err != nil {
				log.Fatal(err)
			}

			// Llamar recursivamente a la función
			findsDir(subDir)
		} else {
			// Si no es un directorio, imprimir el nombre
			fmt.Println(dir.Name() + "/" + fi.Name())
		}
	}
}
