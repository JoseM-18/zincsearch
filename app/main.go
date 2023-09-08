package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var emails = make(chan string, 500)
var dataToZinc = make(chan Email, 500)

// Email is a struct that stores the information extracted from an email message.
type Email struct {
	Date    string `json:"date"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

var waG sync.WaitGroup

// create a struct to store the data

func main() {
	//start profiling
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var rootDirPath string
	flag.StringVar(&rootDirPath, "rootDir", "./allen-p", "path to the root directory")
	flag.Parse()

	createIndex()

	// Start the goroutines to process the emails and insert the data into the search engine
	go findsDir(rootDirPath, emails)
	go processEmails(emails, dataToZinc)
	//
	formatData(dataToZinc)

	waG.Wait()
}

func findsDir(dir string, emails chan string) {
	// Recursively searches for files in a given directory and its subdirectories, and sends the file paths to a channel for further processing.

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
func processEmails(emails chan string, dataToZinc chan Email) {

	
	for email := range emails {
		emailData,err := parseEmail(email)
		if err != nil {
			log.Fatal(err)
			continue
		}
		dataToZinc <- emailData

	}

}

func parseEmail(email string)  (Email, error) {
	/*
		parseEmail parses an email message and extracts relevant information such as the message ID, date, sender, recipient, subject, and body.

		Parameters:
		- email (string): The email message to be parsed.

		Returns:
		- Email: Parsed Email struct containing the extracted information from the email message.
	*/

	fileInfo, err := os.ReadFile(email)
	if err != nil {
		return Email{}, err
	}

	msg, err := mail.ReadMessage(bytes.NewReader(fileInfo))
	if err != nil {
		return Email{}, err
	}

	// Parse the message
	header := msg.Header
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return Email{}, err
	}

	// Create a struct to store the data
	emailData := Email{
		Date:    header.Get("Date"),
		From:    header.Get("From"),
		To:      header.Get("To"),
		Subject: header.Get("Subject"),
		Body:    string(body),
	}

	return emailData, nil
}

func createIndex() {

	/*
		createIndex is responsible for creating an index in a search engine.
		It sends an HTTP POST request to the search engine's API with the necessary information to create the index.

		Example Usage:
		createIndex()

		Inputs:
		None

		Outputs:
		None
	*/

	structureIndex := `{
		"name": "email",
		"storage_type": "disk",
		"shard_num": 1,
		"mappings": {
			"properties": {
				"Date": {
					"type": "date",
					"index": true,
					"sortable": true,
					"aggregatable": true
				},
				"From": {
					"type": "text",
					"index": true,
					"sortable": true,
					"aggregatable": true
				},
				"To": {
					"type": "text",
					"index": true,
					"sortable": true,
					"aggregatable": true
				},
				"Subject": {
					"type": "text",
					"index": true,
					"sortable": true,
					"aggregatable": true
				},
				"Body": {
					"type": "text",
					"index": true,
					"sortable": true,
					"aggregatable": true
				}
			}
		}
	}`
	url := "http://zincsearch:4080/api/index"

	req, err := http.NewRequest("POST", url, strings.NewReader(structureIndex))
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	respuesta, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer respuesta.Body.Close()
}

func formatData(dataToZinc chan Email) {
	/*
		formatData formats the data from an Email struct into a JSON string.

		Parameters:
		- data (Email): The data to be formatted.

		Returns:
		- string: The formatted data as a JSON string.
		- error: An error if there is one, or nil if there is no error.
	*/

	// Get the index information
	// Create a struct to store the data

	// Format the date
	for data := range dataToZinc {

		parsedDate, err := mail.ParseDate(data.Date)
		if err != nil {
			log.Fatal(err)
		}

		date, err := formatDate(parsedDate.String())
		if err != nil {
			log.Fatal(err)
		}

		// Format the data
		formattedData := fmt.Sprintf(`{
		"date": "%s",
		"from": "%s",
		"to": "%s",
		"subject": "%s",
		"body": "%s"

	}`, date, data.From, data.To, data.Subject, data.Body)

		// Send the data to the 'insertData' function
		insertData(formattedData)
	}

}

/**
 * insertData sends an HTTP POST request to the search engine's API to insert data into the index.
 * @param {string} data - The data to be inserted into the index.
 * @returns {void}
 */
func insertData(data string) {

	url := "http://zincsearch:4080/api/email/_doc"

	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		panic(err)
	}

	request.SetBasicAuth("admin", "Complexpass#123")
	request.Header.Set("Content-Type", "application/json")

	respuesta, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer respuesta.Body.Close()

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
