package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"

	"licenseplate.wtf/db"
)

type server struct {
	config *config
	router *mux.Router
}

func NewServer() *server {
	s := server{}
	s.config = loadConfig()
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ListenAndServe() {
	s.routes()
	log.Println("Starting server on ", s.config.Port)
	s.listenWithCancel()
}

// Run the server, asking it to shut down when the interrupt signal is received.
func (s *server) listenWithCancel() {

	// Create a shell http.Server which uses our server to resolve requests.
	srv := http.Server{
		Addr:    s.config.Port,
		Handler: s,
	}

	// Empty channel.
	closed := make(chan struct{})
	go func() {

		// Create a signal channel
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		// Block until the signal is received.
		<-c
		fmt.Printf("\r")
		log.Println("Received interrupt. Shutting down server.")

		// Backup the DB before shutting down in production.
		if os.Getenv("ENV") == "production" {
			db.Backup()
		}

		// Ask the server to shutdown.
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}

		// Once it has, notify the main thread.
		close(closed)
	}()

	// Run the server.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server shutdown error: %v", err)
	}

	// Once this is closed, we know it has been completely shut down.
	<-closed
	log.Println("Server shut down gracefully.")
}
