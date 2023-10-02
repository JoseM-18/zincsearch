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
func FindsDir(dir string, dirsEmails chan string, wgFinders *sync.WaitGroup) {
	defer wgFinders.Done()

	//get the files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		contadorErroresFinders(err)
		return
	}

	//iterate over the files
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Println(err)
			contadorErroresFinders(err)
			return
		}
		//get the file path
		filePath := filepath.Join(dir, file.Name())
		if fileInfo.IsDir() {
			wgFinders.Add(1)
			go FindsDir(filePath, dirsEmails, wgFinders)
		} else {
			//send the file path to the channel
			dirsEmails <- filePath
		}
	}

}

/**
 * contadorErroresFinders stores the errors that occur when searching for email messages.
 * @param {error} err - The error that occurred.
 * @returns {void}
 */
var erros []error

func contadorErroresFinders(err error) {
	erros = append(erros, err)
}

/**
 * GetErroresFinders returns the number of errors that occurred when searching for email messages.
 * @returns {int} - The number of errors.
 */
func GetErroresFinders() int {
	total := len(erros)
	return total
}
