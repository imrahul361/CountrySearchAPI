package server

import (
	"CountrySearch/internal/externalapi"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()

	// Wrap all routes with CORS middleware
	corsWrapper := s.corsMiddleware(r)
	r.HandlerFunc(http.MethodGet, "/api/countries/search", s.SearchCountryHandler)

	return corsWrapper
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Use "*" for all origins, or replace with specific origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are needed

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) SearchCountryHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	var resp externalapi.CountrySearchResponse
	if value, ok := s.cache.Get(name); ok {
		if country, ok := value.(externalapi.CountrySearchResponse); ok {
			resp = country
		}
	} else {
		country, err := externalapi.FetchCountryData(name)
		if err != nil {
			log.Printf("error fetching country data: %v", err)
			http.Error(w, "Country not found", http.StatusNotFound)
			return
		}
		resp = country
		s.cache.Set(name, country)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}
	_, _ = w.Write(jsonResp)
}
