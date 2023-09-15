package apizinc

import (
	"fmt"
	"net/http"
	"strings"
)

/**
 * existIndex sends an HTTP HEAD request to the search engine's API to check if an index exists.
 * @returns {boolean} - Returns true if the index exists, false otherwise.
 */

func existIndex() bool {
	req, err := http.NewRequest("HEAD", "http://localhost:4080/api/index/email", nil)

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

	exist := respuesta.StatusCode == 200
	defer respuesta.Body.Close()

	return exist
}

/**
 * createIndex sends an HTTP POST request to the search engine's API to create an index.
 * @returns {void}
 */
func CreateIndex() {

	if !existIndex() {

		structureIndex := `{
		"name": "email",
		"storage_type": "disk",
		"shard_num": 1,
		"mappings": {
			"properties": {
				"date": {
					"type": "date",
					"index": true,
					"store": false,
				},
				"from": {
					"type": "text",
					"index": true,
					"store": false,
				},
				"to": {
					"type": "text",
					"index": true,
					"store": false,
				},
				"subject": {
					"type": "text",
					"index": true,
					"highlightable": true,
				},
				"body": {
					"type": "text",
					"index": true,
					"store": false,
					"highlightable": true,
				}
			}
		}
	}`
		url := "http://localhost:4080/api/index"

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
}

/**
 * insertData sends an HTTP POST request to the search engine's API to insert data into the index.
 * @param {string} data - The data to be inserted into the index.
 * @returns {void}
 */
func InsertData(data string) {

	url := "http://localhost:4080/api/email/_multi"

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

	fmt.Println(respuesta)

	defer respuesta.Body.Close()

}

func Search(query string) (string,error) {
	structureSearch := `{
		"search_type": "match",
		"query":
		{
				"term": "` + query + `"
		},
		"highlight": {
				"fields": {
						"body": {},
						"subject": {}
				}
		}
}`

	url := "http://localhost:4080/api/email/_search"

	request, err := http.NewRequest("GET", url, strings.NewReader(structureSearch))
	if err != nil {
		return "",err
	}

	request.SetBasicAuth("admin", "Complexpass#123")
	request.Header.Set("Content-Type", "application/json")

	respuesta, err := http.DefaultClient.Do(request)
	if err != nil {
		return "",err
	}

	defer respuesta.Body.Close()

	res := make([]byte, respuesta.ContentLength)
	fmt.Println(respuesta.Body.Read(res))

	return string(res),nil


}
