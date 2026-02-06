package server

import (
	"CountrySearch/internal/cache"
	"CountrySearch/internal/externalapi"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer() *Server {
	return &Server{
		port:  8080,
		cache: cache.NewLRUCache(100),
	}
}

func TestRegisterRoutes_ReturnsHandler(t *testing.T) {
	s := setupTestServer()
	handler := s.RegisterRoutes()
	assert.NotNil(t, handler)
}

func TestCorsMiddleware_SetsHeaders(t *testing.T) {
	s := setupTestServer()
	handler := s.RegisterRoutes()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCorsMiddleware_HandlesOptions(t *testing.T) {
	s := setupTestServer()
	handler := s.RegisterRoutes()

	req := httptest.NewRequest("OPTIONS", "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestSearchCountryHandler_ReturnsCachedData(t *testing.T) {
	s := setupTestServer()

	testData := externalapi.CountrySearchResponse{
		Name:    "France",
		Capital: "Paris",
	}
	s.cache.Set("france", testData)

	req := httptest.NewRequest("GET", "/api/countries/search?name=france", nil)
	rr := httptest.NewRecorder()
	s.SearchCountryHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response externalapi.CountrySearchResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "France", response.Name)
}

func TestSearchCountryHandler_ParsesJSONResponse(t *testing.T) {
	s := setupTestServer()

	testData := externalapi.CountrySearchResponse{
		Name: "Germany",
	}
	s.cache.Set("germany", testData)

	req := httptest.NewRequest("GET", "/api/countries/search?name=germany", nil)
	rr := httptest.NewRecorder()
	s.SearchCountryHandler(rr, req)

	var response externalapi.CountrySearchResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCache_StoresAndRetrievesData(t *testing.T) {
	s := setupTestServer()

	testData := externalapi.CountrySearchResponse{
		Name: "Spain",
	}
	s.cache.Set("spain", testData)

	cached, ok := s.cache.Get("spain")
	assert.True(t, ok)
	assert.NotNil(t, cached)
}

func TestSearchCountryHandler_WithMultipleCountries(t *testing.T) {
	s := setupTestServer()

	countries := []externalapi.CountrySearchResponse{
		{Name: "USA", Capital: "Washington"},
		{Name: "UK", Capital: "London"},
		{Name: "Japan", Capital: "Tokyo"},
	}

	for _, country := range countries {
		s.cache.Set(country.Name, country)
	}

	for _, country := range countries {
		cached, ok := s.cache.Get(country.Name)
		assert.True(t, ok)
		assert.NotNil(t, cached)
	}
}

func TestCorsMiddleware_AllowsPostRequest(t *testing.T) {
	s := setupTestServer()
	handler := s.RegisterRoutes()

	req := httptest.NewRequest("POST", "/api/countries/search", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestSearchCountryHandler_ResponseBodyNotEmpty(t *testing.T) {
	s := setupTestServer()

	s.cache.Set("italy", externalapi.CountrySearchResponse{
		Name: "Italy",
	})

	req := httptest.NewRequest("GET", "/api/countries/search?name=italy", nil)
	rr := httptest.NewRecorder()
	s.SearchCountryHandler(rr, req)

	assert.NotEmpty(t, rr.Body.String())
}
