package email

import (
	"io"
	"net/mail"
	"os"
)

type Email struct {
	MessageId string `json:"messageId"`
	Date      string `json:"date"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

/**
 * parseEmail parses an email message and extracts the relevant information from it.
 * @param {string} email - The path to the email message.
 * @returns {Email} - The extracted information from the email message.
 * @returns {error} - An error if there is one, or nil if there is no error.
 */
func ParseEmail(email string) (Email, error) {

	// Open the file
	file, err := os.Open(email)
	if err != nil {
		contadorErroresEmail(err)
		return Email{}, err
	}
	defer file.Close()

	// Read the message
	msg, err := mail.ReadMessage(file)
	if err != nil {
		contadorErroresEmail(err)
		return Email{}, err
	}

	// Parse the message
	header := msg.Header
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		contadorErroresEmail(err)
		return Email{}, err
	}

	// Create a struct to store the data
	emailData := Email{
		MessageId: header.Get("Message-ID"),
		Date:      header.Get("Date"),
		From:      header.Get("From"),
		To:        header.Get("To"),
		Subject:   header.Get("Subject"),
		Body:      string(body),
	}

	// Return the data
	return emailData, nil

}

var errors []error

func contadorErroresEmail(err error) {
	errors = append(errors, err)
}

func GetErroresEmail() int {
	total := len(errors)
	return total
}
