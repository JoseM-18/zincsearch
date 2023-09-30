package routes

import (
	"bytes"
	"encoding/json"
	"github.com/JoseM-18/zincSearch/apiZinc"
	"github.com/go-chi/chi/v5"
	"net/http"
	"github.com/go-chi/cors"

)

/**
 * SetupRouter creates a new router for the application.
 * @returns {chi.Mux} - Returns a new router.
 */

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	// Configura CORS
	cors := cors.New(cors.Options{
		// Ajusta los valores según tus necesidades
		AllowedOrigins:   []string{"*"},                                       // Permitir todas las origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Permite los métodos que desees
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Máximo de tiempo en caché de la respuesta de pre-verificación CORS en segundos
	})

	// Usa el middleware CORS en tu router
	router.Use(cors.Handler)

	router.Get("/search", SearchHandler)
	return router
}

/**
 * SearchHandler handles the search request.
 * @param {http.ResponseWriter} w - The response writer.
 * @param {http.Request} r - The request.
 * @returns {void}
 */
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := apizinc.Search(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultsJson, err := printResultsJson(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resultsJson))

}

func printResultsJson(results map[string]interface{}) (string, error) {
	// Create a buffer to store the JSON output
	var buffer bytes.Buffer

	// Marshal the results to JSON
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(results)
	if err != nil {
		return "", err
	}

	// Return a string with the JSON results
	return buffer.String(), nil
}
