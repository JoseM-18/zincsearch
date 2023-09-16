package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"
	"github.com/JoseM-18/zincSearch/apiZinc"
	"github.com/JoseM-18/zincSearch/email"
)

var emails = make(chan string, 100)
var dataToZinc = make(chan email.Email, 20)

// create a struct to store the data

func main() {
	//start profiling
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var rootDirPath string
	flag.StringVar(&rootDirPath, "rootDir", "../enron_mail_20110402/enron_mail_20110402/maildir", "path to the root directory")
	flag.Parse()

	apizinc.CreateIndex()

	// Create a WaitGroup
	var wg, wgfd, wgpe sync.WaitGroup

	// Start the goroutines

	wg.Add(1)
	go formatData(&wg, dataToZinc)

	for i := 0; i < 10; i++ {
		wgpe.Add(1)
		go processEmails(&wgpe, emails, dataToZinc)
	}

	wgfd.Add(1)
	go findsDir(rootDirPath, emails, &wgfd)

	wgfd.Wait()
	close(emails)

	wgpe.Wait()
	close(dataToZinc)

	wg.Wait()
}

/**
 * findsDir recursively searches for files in a given directory and its subdirectories, and sends the file paths to a channel for further processing.
 * @param {string} dir - The directory to be searched.
 * @param {chan string} emails - A channel containing the paths to the email messages.
 * @returns {void}
 */
func findsDir(dir string, emails chan string, wgfd *sync.WaitGroup) {
	defer wgfd.Done()

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Println(err)
			return
		}

		if fileInfo.IsDir() {
			wgfd.Add(1)
			go findsDir(filepath.Join(dir, file.Name()), emails, wgfd)
		} else {
			emails <- filepath.Join(dir, file.Name())
		}
	}

}

/**
 * processEmails processes the emails and extracts the relevant information from them.
 * @param {chan string} emails - A channel containing the paths to the email messages.
 * @param {chan Email} dataToZinc - A channel containing the extracted information from the email messages.
 * @returns {void}
 */
func processEmails(wgpe *sync.WaitGroup, emails chan string, dataToZinc chan email.Email) {
	defer wgpe.Done()
	for oneEmail := range emails {
		emailData, err := email.ParseEmail(oneEmail)
		if err != nil {
			log.Println(err)
		} else {
			dataToZinc <- emailData
		}
	}

}

/**
 * formatData formats the data and sends it to the search engine.
 * @param {chan Email} dataToZinc - A channel containing the extracted information from the email messages.
 * @returns {void}
 */
func formatData(wg *sync.WaitGroup, dataToZinc chan email.Email) {
	defer wg.Done()

	buffer := make([]email.Email, 0)
	for data := range dataToZinc {
		buffer = append(buffer, data)
		if len(buffer) == 1000 {
			sendBufferedData(buffer)
			buffer = buffer[:0]
		}
	}

	if len(buffer) > 0 {
		sendBufferedData(buffer)

	}

}

/**
 * sendBufferedData sends the data to the 'insertData' function.
 * @param {[]Email} buffer - A slice containing the extracted information from the email messages.
 * @returns {void}
 */
func sendBufferedData(dataBuffer []email.Email) {
	var buffer bytes.Buffer
	for _, item := range dataBuffer {

		jsonData, err := json.Marshal(item)
		if err != nil {
			log.Fatal(err)
		}
		buffer.Write(jsonData)
		buffer.WriteString("\n") // Add a newline after each JSON object because the search engine expects it
	}

	// Send the data to the search engine for indexing and searching
	apizinc.InsertData(buffer.String())
}
