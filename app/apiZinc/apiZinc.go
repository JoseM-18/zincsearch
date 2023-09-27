package apizinc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const host = "zincsearch"
const port = "4080"

/**
 * existIndex sends an HTTP HEAD request to the search engine's API to check if an index exists.
 * @returns {boolean} - Returns true if the index exists, false otherwise.
 */
func existIndex() (bool,error) {
	
	req, err := http.NewRequest("HEAD", "http://"+host+":"+port+"/api/index/email", nil)

	if err != nil {
		log.Println(err)
		contadorErroresApiZinc(err)
		return false, err
	}

	respuesta, err := requestZinc(req)
	if err != nil {
		log.Println(err)
		contadorErroresApiZinc(err)
		return false, err
	}

	defer respuesta.Body.Close()

	return respuesta.StatusCode == 200, nil
}

/**
 * createIndex sends an HTTP POST request to the search engine's API to create an index.
 * @returns {void}
 */
func CreateIndex() (string, error) {

	value,err := existIndex()
	if err != nil {
		log.Println(err)
		contadorErroresApiZinc(err)
		return "", err
	}

	if !value { 

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
		url := "http://" + host + ":4080/api/index"

		req, err := http.NewRequest("POST", url, strings.NewReader(structureIndex))
		if err != nil {
			log.Println(err)
			contadorErroresApiZinc(err)
			return "", err
		}

		respuesta, err := requestZinc(req)
		if err != nil {
			log.Println(err)
			contadorErroresApiZinc(err)
			return "", err
		}
		defer respuesta.Body.Close()

		return "Index created", nil

	} else {
		return "Index already exists", nil
	}
}
/**
 * insertData sends an HTTP POST request to the search engine's API to insert data into the index.
 * @param {string} data - The data to be inserted into the index.
 * @returns {void}
 */
func InsertData(data string) error {

	const url = "http://" + host + ":4080/api/email/_multi"

	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		contadorErroresApiZinc(err)
		return fmt.Errorf("error in request: %w", err)
	}

	respuesta, err := requestZinc(request)
	if err != nil {
		contadorErroresApiZinc(err)
		return fmt.Errorf("error in response: %w", err)
	}

	defer respuesta.Body.Close()

	return nil
}

/**
 * search sends an HTTP POST request to the search engine's API to search for a query. 
 * @param {string} query - The query to be searched.
 * @returns {map[string]interface} - The results of the search.
 * @returns {error} - The error.
 */
func Search(query string) (map[string]interface{}, error) {
	structureSearch := `{
		"search_type":"match",
		"query":{
			 "term":"` + query + `"
		},
		"max_results":1000,
		"highlight":{
			 "fields":{
				 "from":{
						
				 },
				 "to":{
						
				 },
				"body":{
						 
				}
			 }
		}
 }`

	url := "http://" + host + ":4080/api/email/_search"

	request, err := http.NewRequest("POST", url, strings.NewReader(structureSearch))
	if err != nil {
		contadorErroresApiZinc(err)
		return nil, err
	}

	respuesta, err := requestZinc(request)
	if err != nil {
		contadorErroresApiZinc(err)
		return nil, err
	}
	defer respuesta.Body.Close()

	// Decode the response into a map of strings and interfaces
	var results map[string]interface{}
	err = json.NewDecoder(respuesta.Body).Decode(&results)
	if err != nil {
		contadorErroresApiZinc(err)
		return nil, err
	}

	return results, nil

}



/**
 * requestZinc sends an HTTP request to the search engine's API.
 * @param {http.Request} request - The HTTP request to be sent.
 * @returns {http.Response} - The HTTP response.
 * @returns {error} - The error.
 */
func requestZinc (resquest *http.Request) (*http.Response, error) {
	username := os.Getenv("ZINCSEARCH_USERNAME")
	password := os.Getenv("ZINCSEARCH_PASSWORD")
	if username == "" || password == "" {
	    return nil, fmt.Errorf("ZINCSEARCH_USERNAME or ZINCSEARCH_PASSWORD environment variables are not set or empty")
	}
	resquest.SetBasicAuth(username, password)
	resquest.Header.Set("Content-Type", "application/json")
	resquest.Close = true
	
	return http.DefaultClient.Do(resquest)
}

var errors []error
func contadorErroresApiZinc(err error){
	errors = append(errors, err)
}

func GetErroresApiZinc() int{
	total := len(errors)
	return total
}

