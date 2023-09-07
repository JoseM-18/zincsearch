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
	"strings"
	"sync"
	"time"

)

var emails = make(chan string, 500)
var dataToZinc = make(chan Email, 500)

var waG sync.WaitGroup

// create a struct to store the data
type Email struct {
	ID        string
	MessageId string
	Date      string
	From      string
	To        string
	Subject   string
	Body      string
}

func main() {
	//start profiling
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var rootDirPath string
	flag.StringVar(&rootDirPath, "rootDir", "./unarchivo", "path to the root directory")
	flag.Parse()

	rootDir, err := os.Open(rootDirPath)
	if err != nil {
		log.Println(err)
		return
	}

	findsDir(rootDir)
	// create index
	createIndex()

	
	go processEmails(emails)

	go insertData(dataToZinc)

	os.Stdout.Sync()

	waG.Wait()
}

func findsDir(dir *os.File) {
	// Recursively searches for files in a given directory and its subdirectories, and sends the file paths to a channel for further processing.

	// Inputs:
	// - dir: A pointer to an os.File representing the directory to search.

	// Outputs:
	// - None.
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	_ = dir.Close()

	total := 0
	// Iterate over the files in the directory
	for _, file := range files {
		// process each file or directory
		if file.IsDir() {
			// if the file is a directory, open it and call the function recursively
			subDir, err := os.Open(dir.Name() + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}

			findsDir(subDir)
		} else {

			total++
			waG.Add(1)
			emails <- dir.Name() + "/" + file.Name()

		}
	}

}

func processEmails(emails chan string) {
	// processEmails processes emails from the 'emails' channel.
	// It reads each email from the channel, parses it using the 'parseEmail' function,
	// and sends the parsed email data to the 'dataToZinc' channel.
	// The function also keeps track of the number of processed emails using the 'times' variable.

	for email := range emails {
		dataToZinc <- parseEmail(email)
		waG.Done()
	}
	fmt.Println("hahahah",len(dataToZinc))
}

func parseEmail(email string) Email {
	/*
		parseEmail parses an email message and extracts relevant information such as the message ID, date, sender, recipient, subject, and body.

		Parameters:
		- email (string): The email message to be parsed.

		Returns:
		- Email: Parsed Email struct containing the extracted information from the email message.
	*/

	fileInfo, err := os.ReadFile(email)
	if err != nil {
		log.Fatal(err)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(fileInfo))
	if err != nil {
		log.Fatal(err)
	}

	// Parse the message
	header := msg.Header
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a struct to store the data
	emailData := Email{
		Date:    header.Get("Date"),
		From:    header.Get("From"),
		To:      header.Get("To"),
		Subject: header.Get("Subject"),
		Body:    string(body),
	}

	return emailData
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

func insertData(dataToZinc chan Email) {
	defer waG.Done()
	// Send an HTTP POST request to insert data into an index in a search engine.
	//
	// Parameters:
	//   - data (string): The data to be inserted into the index.
	//
	// Example Usage:
	//   insertData(data)
	//
	// Flow:
	//   1. Create an HTTP POST request with the specified URL and data.
	//   2. Set the basic authentication credentials and headers for the request.
	//   3. Send the request using the default HTTP client.
	//   4. Close the response body after the request is completed.
	// Get the index information
	// Create a struct to store the data
		data := <-dataToZinc
		parsedDate, err := mail.ParseDate(data.Date)
		if err != nil {
			panic(err)
		}
		formattedDate, err := formatDate(parsedDate.String())
		if err != nil {
			panic(err)
		}
		dataJson := fmt.Sprintf(`{"date": "%s", "from": "%s", "to": "%s", "subject": "%s"}`, formattedDate, data.From, data.To, data.Subject)
		fmt.Println(dataJson)
		url := "http://zincsearch:4080/api/email/_doc"

		request, err := http.NewRequest("POST", url, strings.NewReader(dataJson))
		if err != nil {
			panic(err)
		}

		request.SetBasicAuth("admin", "Complexpass#123")
		request.Header.Set("Content-Type", "application/json")

		respuesta, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
		fmt.Println(respuesta)
		
}

func formatDate(date string) (time.Time, error) {
	const customDateFormat = "2006-01-02 15:04:05 -0700 -0700"
	t, err := time.ParseInLocation(customDateFormat, date, time.UTC )
	if err != nil {
		return time.Time{}, err
	}
	formattedDate := t.Format(time.RFC3339)
	fmt.Println(formattedDate)
	return t, nil
}
