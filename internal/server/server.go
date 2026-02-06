package server

import (
	"CountrySearch/internal/cache"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port  int
	cache *cache.LRUCache
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}
	cache := cache.NewLRUCache(100)
	NewServer := &Server{
		port:  port,
		cache: cache,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
