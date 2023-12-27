package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
    db *sql.DB
}

func NewServer(port int, db *sql.DB) *http.Server {
	NewServer := &Server{
		port: port,
        db: db,
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
