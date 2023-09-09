package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/JoseM-18/zincSearch/apiZinc"
	"github.com/JoseM-18/zincSearch/email"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/mail"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var emails = make(chan string, 500)
var dataToZinc = make(chan email.Email, 500)

var waG sync.WaitGroup

// create a struct to store the data

func main() {
	//start profiling
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var rootDirPath string
	flag.StringVar(&rootDirPath, "rootDir", "../asd", "path to the root directory")
	flag.Parse()

	apizinc.CreateIndex()

	// Start the goroutines to process the emails and insert the data into the search engine
	go findsDir(rootDirPath, emails)
	go processEmails(emails, dataToZinc)
	//
	formatData(dataToZinc)

	waG.Wait()
}

/**
 * findsDir recursively searches for files in a given directory and its subdirectories, and sends the file paths to a channel for further processing.
 * @param {string} dir - The directory to be searched.
 * @param {chan string} emails - A channel containing the paths to the email messages.
 * @returns {void}
 */
func findsDir(dir string, emails chan string) {

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
			findsDir(filepath.Join(dir, file.Name()), emails)
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
func processEmails(emails chan string, dataToZinc chan email.Email) {

	for oneEmail := range emails {
		emailData, err := email.ParseEmail(oneEmail)
		if err != nil {
			log.Fatal(err)
			continue
		}
		dataToZinc <- emailData

	}

}

/**
 * formatData formats the data and sends it to the 'insertData' function.
 * @param {chan Email} dataToZinc - A channel containing the extracted information from the email messages.
 * @returns {void}
 */
func formatData(dataToZinc chan email.Email) {

	for data := range dataToZinc {

		// Parse the date
		parsedDate, err := mail.ParseDate(data.Date)
		if err != nil {
			log.Fatal(err)
		}

		// Modify the date
		date, err := formatDate(parsedDate.String())
		if err != nil {
			log.Fatal(err)
		}

		// Set the modified date
		data.Date = date.String()

		// JSON marshal the email object
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		// Send the data to the 'insertData' function
		apizinc.InsertData(string(jsonData))

	}

}

/**
 * formatDate formats a date string into a RFC3339 format.
 * @param {string} date - The date string to be formatted.
 * @returns {time.Time} - The formatted date.
 * @returns {error} - An error if there is one, or nil if there is no error.
 */
func formatDate(date string) (time.Time, error) {
	const customDateFormat = "2006-01-02 15:04:05 -0700 -0700"
	t, err := time.ParseInLocation(customDateFormat, date, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	formattedDate := t.Format(time.RFC3339)
	fmt.Println(formattedDate)
	return t, nil
}
