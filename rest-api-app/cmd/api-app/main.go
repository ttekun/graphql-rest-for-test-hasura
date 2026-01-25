package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ttekun/rest-api-app/internal/server"
)

func main() {
	// Create server instance
	srv := server.NewServer("8080")

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started successfully")

	// Wait for interrupt signal
	<-done
	log.Println("Server stopping...")
}