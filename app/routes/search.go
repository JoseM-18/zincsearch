package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/JoseM-18/zincSearch/apiZinc"
	"net/http"
	"encoding/json"
)

/**
 * SetupRouter creates a new router for the application.
 * @returns {chi.Mux} - Returns a new router.
 */
 
func SetupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/search", SearchHandler)
	return router
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results,err := apizinc.Search(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}




