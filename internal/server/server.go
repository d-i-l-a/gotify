package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gotify/internal/cache"
	"gotify/internal/database"
)

type Server struct {
	port  int
	db    database.Service
	cache cache.Cache
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	redisAddr := os.Getenv("REDISURL")
	newServer := &Server{
		port:  port,
		db:    database.New(),
		cache: cache.New(redisAddr),
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
