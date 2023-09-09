package apizinc

import (
	"net/http"
	"strings"
)

/**
 * createIndex sends an HTTP POST request to the search engine's API to create an index.
 * @returns {void}
 */
func CreateIndex() {

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

/**
 * insertData sends an HTTP POST request to the search engine's API to insert data into the index.
 * @param {string} data - The data to be inserted into the index.
 * @returns {void}
 */
func InsertData(data string) {

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
