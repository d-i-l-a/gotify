package server

import (
	"fmt"
	"gotify/internal/cache"
	"gotify/internal/database"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port int
	db   database.Service
	c    cache.Cache
}

func NewServer() *http.Server {
	cache := cache.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port: port,
		db:   database.New(),
		c:    cache,
	}
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
