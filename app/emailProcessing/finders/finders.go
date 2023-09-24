package finders

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

/**
 * findsDir recursively searches for files in a given directory and its subdirectories, and sends the file paths to a channel for further processing.
 * @param {string} dir - The directory to be searched.
 * @param {chan string} emails - A channel containing the paths to the email messages.
 * @returns {void}
 */
func FindsDir(dir string, emails chan string, wgFinders *sync.WaitGroup) {
	defer wgFinders.Done()

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		contadorErroresFinders(err)
		return
	}

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Println(err)
			contadorErroresFinders(err)
			return
		}

		if fileInfo.IsDir() {
			wgFinders.Add(1)
			go FindsDir(filepath.Join(dir, file.Name()), emails, wgFinders)
		} else {
			emails <- filepath.Join(dir, file.Name())
		}
	}

}

var erros []error
func contadorErroresFinders(err error){
	erros = append(erros, err)
}

func GetErroresFinders() int{
	total := len(erros)
	return total
}