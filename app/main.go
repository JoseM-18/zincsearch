package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/mail"
	"os"
	"sync"
)

var emails = make(chan string, 500)
var dataToZinc = make(chan Email, 500)

var waG sync.WaitGroup

// create a struct to store the data
type Email struct {
	ID       string
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

	rootDir, err := os.Open("./allen-p")
	if err != nil {
		log.Fatal(err)
	}

	go processEmails()

	findsDir(rootDir)
	close(emails)

	os.Stdout.Sync()

	// loop infinitely
	select {}
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

func processEmails() {
	// processEmails processes emails from the 'emails' channel.
	// It reads each email from the channel, parses it using the 'parseEmail' function,
	// and sends the parsed email data to the 'dataToZinc' channel.
	// The function also keeps track of the number of processed emails using the 'times' variable.
	times := 0

	for email := range emails {
		times++
		waG.Done()
		dataToZinc <- parseEmail(email)
	}
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
		ID:       email,
		MessageId: header.Get("Message-Id"),
		Date:      header.Get("Date"),
		From:      header.Get("From"),
		To:        header.Get("To"),
		Subject:   header.Get("Subject"),
		Body:      string(body),
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
	return emailData
}
