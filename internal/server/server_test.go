package server

import (
	"CountrySearch/internal/cache"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer_CreatesHTTPServer(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.Addr)
}

func TestNewServer_SetsTimeouts(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.Equal(t, 10*time.Second, server.ReadTimeout)
	assert.Equal(t, 30*time.Second, server.WriteTimeout)
	assert.Equal(t, time.Minute, server.IdleTimeout)
}

func TestNewServer_HasHandler(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.NotNil(t, server.Handler)
}

func TestNewServer_WithoutPort(t *testing.T) {
	os.Unsetenv("PORT")
	server := NewServer()
	assert.Equal(t, ":0", server.Addr)
}

func TestNewServer_WithPort3000(t *testing.T) {
	os.Setenv("PORT", "3000")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.Equal(t, ":3000", server.Addr)
}

func TestRegisterRoutes_ReturnsRouter(t *testing.T) {
	s := &Server{
		port:  8080,
		cache: cache.NewLRUCache(100),
	}

	router := s.RegisterRoutes()
	assert.NotNil(t, router)
}

func TestNewServer_ValidHTTPServer(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.IsType(t, &http.Server{}, server)
}

func TestNewServer_Port8000(t *testing.T) {
	os.Setenv("PORT", "8000")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.Equal(t, ":8000", server.Addr)
}

func TestServer_ReadTimeoutIs10Seconds(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.Equal(t, 10*time.Second, server.ReadTimeout)
}

func TestServer_WriteTimeoutIs30Seconds(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	server := NewServer()
	assert.Equal(t, 30*time.Second, server.WriteTimeout)
}

func TestNewServer_MultipleInstances(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	server1 := NewServer()
	server2 := NewServer()

	require.NotNil(t, server1)
	require.NotNil(t, server2)
	assert.Equal(t, server1.Addr, server2.Addr)
}
