package email

import (
	"bytes"
	"io"
	"net/mail"
	"os"
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

	return emailData, nil
}
