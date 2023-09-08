
package api

import (
	"net/http"
	"strings"
)


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