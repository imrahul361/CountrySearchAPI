# CountrySearchAPI
REST API service that provides country information

# Run the Server
go run ./cmd/api/main.go

# Request
curl -X GET http://localhost:8080/api/countries/search?name=India
