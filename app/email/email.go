package email

import (
	"bytes"
	"io"
	"log"
	"net/mail"
	"os"
	"time"
)

type Email struct {
	Date    string `json:"date"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

/**
 * parseEmail parses an email message and extracts the relevant information from it.
 * @param {string} email - The path to the email message.
 * @returns {Email} - The extracted information from the email message.
 * @returns {error} - An error if there is one, or nil if there is no error.
 */
func ParseEmail(email string) (Email, error) {

	// Read the file
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

	// Format the date
	emailData, err = formatDate(emailData)
	if err != nil {
		return Email{}, err
	}

	return emailData, nil

}

/**
 * formatDate formats a date string into a RFC3339 format.
 * @param {string} date - The date string to be formatted.
 * @returns {time.Time} - The formatted date.
 * @returns {error} - An error if there is one, or nil if there is no error.
 */
func formatDate(data Email) (Email, error) {
	// Parse the date
	parsedDate, err := mail.ParseDate(data.Date)
	if err != nil {
		log.Fatal(err)
	}
	// Modify the date
	date, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", parsedDate.String())
	if err != nil {
		log.Fatal(err)
	}

	// Set the modified date
	data.Date = date.String()

	return data, nil
}
