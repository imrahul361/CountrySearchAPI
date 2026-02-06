package externalapi

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRoundTripper for mocking HTTP requests
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestFetchCountryDataWithClient_Success(t *testing.T) {
	mockResponse := `[
        {
            "name": "United States",
            "capital": "Washington, D.C.",
            "population": 331000000,
            "currencies": [{"symbol": "$"}]
        }
    ]`

	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockResponse)),
				}, nil
			},
		},
	}

	result, err := FetchCountryDataWithClient("United States", client)

	assert.NoError(t, err)
	assert.Equal(t, "United States", result.Name)
	assert.Equal(t, "Washington, D.C.", result.Capital)
	assert.Equal(t, 331000000, result.Population)
	assert.Equal(t, "$", result.Currency)
}

func TestFetchCountryDataWithClient_CountryNotFound(t *testing.T) {
	mockResponse := `[
        {
            "name": "United States",
            "capital": "Washington, D.C.",
            "population": 331000000,
            "currencies": [{"symbol": "$"}]
        }
    ]`

	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockResponse)),
				}, nil
			},
		},
	}

	result, err := FetchCountryDataWithClient("NonExistentCountry", client)

	assert.NoError(t, err)
	assert.Equal(t, "", result.Name)
}

func TestFetchCountryDataWithClient_CaseInsensitive(t *testing.T) {
	mockResponse := `[
        {
            "name": "France",
            "capital": "Paris",
            "population": 67000000,
            "currencies": [{"symbol": "€"}]
        }
    ]`

	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockResponse)),
				}, nil
			},
		},
	}

	result, err := FetchCountryDataWithClient("france", client)

	assert.NoError(t, err)
	assert.Equal(t, "France", result.Name)
	assert.Equal(t, "Paris", result.Capital)
}

func TestFetchCountryDataWithClient_NoCurrencies(t *testing.T) {
	mockResponse := `[
        {
            "name": "TestCountry",
            "capital": "TestCapital",
            "population": 1000000,
            "currencies": []
        }
    ]`

	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockResponse)),
				}, nil
			},
		},
	}

	result, err := FetchCountryDataWithClient("TestCountry", client)

	assert.NoError(t, err)
	assert.Equal(t, "TestCountry", result.Name)
	assert.Equal(t, "", result.Currency)
}

func TestFetchCountryDataWithClient_MultipleCurrencies(t *testing.T) {
	mockResponse := `[
        {
            "name": "MultiCurrency",
            "capital": "Capital",
            "population": 5000000,
            "currencies": [{"symbol": "$"}, {"symbol": "€"}]
        }
    ]`

	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockResponse)),
				}, nil
			},
		},
	}

	result, err := FetchCountryDataWithClient("MultiCurrency", client)

	assert.NoError(t, err)
	assert.Equal(t, "$", result.Currency)
}

func TestFetchCountryDataWithClient_HTTPError(t *testing.T) {
	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
		},
	}

	_, err := FetchCountryDataWithClient("TestCountry", client)

	assert.Error(t, err)
}

func TestFetchCountryDataWithClient_EmptyName(t *testing.T) {
	client := &http.Client{}

	_, err := FetchCountryDataWithClient("", client)

	assert.Error(t, err)
	assert.Equal(t, "country name cannot be empty", err.Error())
}

func TestFetchCountryDataWithClient_InvalidJSON(t *testing.T) {
	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("invalid json")),
				}, nil
			},
		},
	}

	_, err := FetchCountryDataWithClient("TestCountry", client)

	assert.Error(t, err)
}

func TestCountrySearchResponse_Structure(t *testing.T) {
	response := CountrySearchResponse{
		Name:       "Germany",
		Capital:    "Berlin",
		Currency:   "€",
		Population: 83000000,
	}

	assert.Equal(t, "Germany", response.Name)
	assert.Equal(t, "Berlin", response.Capital)
	assert.Equal(t, "€", response.Currency)
	assert.Equal(t, 83000000, response.Population)
}

func TestCountryAPIResponse_Structure(t *testing.T) {
	apiResponse := CountryAPIResponse{
		Name:       "Italy",
		Capital:    "Rome",
		Population: 60000000,
	}

	assert.Equal(t, "Italy", apiResponse.Name)
	assert.Equal(t, "Rome", apiResponse.Capital)
	assert.Equal(t, 60000000, apiResponse.Population)
}

func TestFetchCountryDataWithClient_NetworkError(t *testing.T) {
	client := &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		},
	}

	_, err := FetchCountryDataWithClient("Test", client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch country data")
}
