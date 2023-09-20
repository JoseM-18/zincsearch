package processor

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"github.com/JoseM-18/zincSearch/apiZinc"
	"github.com/JoseM-18/zincSearch/email"
)

/**
 * processEmails processes the emails and extracts the relevant information from them.
 * @param {chan string} emails - A channel containing the paths to the email messages.
 * @param {chan Email} dataToZinc - A channel containing the extracted information from the email messages.
 * @returns {void}
 */
func ProcessEmails(wgProcessors *sync.WaitGroup, emails chan string, dataToZinc chan email.Email) {
	defer wgProcessors.Done()
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
 * sendPackages sends the data to the search engine for indexing and searching.
 * @param {chan Email} dataToZinc - A channel containing the extracted information from the email messages.
 * @returns {void}
 */
func SendPackages(wgSender *sync.WaitGroup, dataToZinc chan email.Email) {
	defer wgSender.Done()
	buffer := make([]email.Email, 0, 10000)

	for data := range dataToZinc {
		buffer = append(buffer, data)
		if len(buffer) == 10000 {
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
