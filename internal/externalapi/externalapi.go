package externalapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CountrySearchResponse struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int    `json:"population"`
}

type CountryAPIResponse struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Population int    `json:"population"`
	Currencies []struct {
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

// FetchCountryDataWithClient allows dependency injection for testing
func FetchCountryDataWithClient(name string, client *http.Client) (CountrySearchResponse, error) {
	if name == "" {
		return CountrySearchResponse{}, fmt.Errorf("country name cannot be empty")
	}

	resp, err := client.Get("https://www.apicountries.com/name/" + name)
	if err != nil {
		return CountrySearchResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CountrySearchResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return CountrySearchResponse{}, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var apiResponse []CountryAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return CountrySearchResponse{}, err
	}

	for _, country := range apiResponse {
		if strings.EqualFold(country.Name, name) {
			response := CountrySearchResponse{
				Name:       country.Name,
				Capital:    country.Capital,
				Population: country.Population,
			}
			if len(country.Currencies) > 0 {
				response.Currency = country.Currencies[0].Symbol
			}
			return response, nil
		}
	}

	return CountrySearchResponse{}, nil
}

// FetchCountryData uses default HTTP client
func FetchCountryData(name string) (CountrySearchResponse, error) {
	return FetchCountryDataWithClient(name, http.DefaultClient)
}
