package processor

import (
	"encoding/json"
	"github.com/JoseM-18/zincSearch/apiZinc"
	"github.com/JoseM-18/zincSearch/email"
	"log"
	"strings"
	"sync"
)

/**
 * processEmails processes the emails and extracts the relevant information from them.
 * @param {chan string} emails - A channel containing the paths to the email messages.
 * @param {chan Email} dataToZinc - A channel containing the extracted information from the email messages.
 * @returns {void}
 */
func ProcessEmails(wgProcessors *sync.WaitGroup, emails chan string, dataToZinc chan email.Email, dateMap *sync.Map) {
	defer wgProcessors.Done()

	// Iterate over the emails
	for oneEmail := range emails {
		// Parse the email
		emailData, err := email.ParseEmail(oneEmail)
		if err != nil {
			log.Println(err)
			contadorErroresProcessor(err)
		} else {
			// Get the date
			date := emailData.Date

			// verify if the date is already in the map and if not, add it
			_, loaded := dateMap.LoadOrStore(date, true)
			if !loaded {
				dataToZinc <- emailData
			}
		}
	}

}

/**
 * sendPackages sends the data to the search engine for indexing and searching.
 * @param {chan Email} dataToZinc - A channel containing the extracted information from the email messages.
 * @returns {void}
 */
func SendPackages(wgSender *sync.WaitGroup, dataToZinc chan email.Email) {
	defer wgSender.Done()
	var buffer []email.Email
	// Send the data in batches of 8000
	const maxBufferSize = 8000

	// Iterate over the data
	for data := range dataToZinc {
		buffer = append(buffer, data)

		// If the buffer is full, send the data to the search engine for indexing and searching
		if len(buffer) == maxBufferSize {
			err := sendBufferedData(buffer)
			if err != nil {
				log.Println(err)
			}
			buffer = []email.Email{}
		}
	}

	// Send the remaining data to the search engine for indexing and searching
	if len(buffer) > 0 {
		err := sendBufferedData(buffer)
		if err != nil {
			log.Println(err)
		}
	}

}

/**
 * sendBufferedData sends the data to the search engine for indexing and searching.
 * @param {[]Email} buffer - A slice containing the extracted information from the email messages.
 * @returns {void}
 */
func sendBufferedData(dataBuffer []email.Email) error {
	var builder strings.Builder

	// Iterate over the data
	for _, item := range dataBuffer {

		// Convert the data to JSON
		jsonData, err := json.Marshal(item)
		if err != nil {
			log.Println(err)
			contadorErroresProcessor(err)
			return err
		}

		// Add the JSON data to the string builder separated by a new line
		builder.Write(jsonData)
		builder.WriteString("\n")
	}

	// Send the data to the search engine for indexing and searching
	err := apizinc.InsertData(builder.String())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

/**
 * contadorErroresProcessor stores the errors that occur when processing email messages.
 * @param {error} err - The error that occurred.
 * @returns {void}
 */
var errors []error

func contadorErroresProcessor(err error) {
	errors = append(errors, err)
}

/**
 * GetErroresProcessor returns the number of errors that occurred when processing email messages.
 * @returns {int} - The number of errors.
 */
func GetErroresProcessor() int {
	total := len(errors)
	return total
}
