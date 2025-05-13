package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener-go/api"
	"url-shortener-go/config"

	"github.com/getsentry/sentry-go"
)

// main is the entry point function for the URL shortening service.
//
// This function performs the following main tasks:
//  1. Initialize Sentry error tracking system
//  2. Create a cancellable context and handle interrupt signals
//  3. Set up deferred functions for closing database and cache connections
//  4. Run the HTTP server
//
// Main features:
//   - Ensures all resources are properly released when the application terminates
//   - Detects interrupt signals (Ctrl+C or SIGTERM) for graceful shutdown
//   - Error monitoring and logging through Sentry
//   - Safe termination of database and cache connections
//
// This function has no return value and runs until the program terminates.
func main() {
	// Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.GetEnv("SENTRY_DSN"),
	}); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
	// Ensure flush before exit
	defer sentry.Flush(2 * time.Second)

	// Create a context that will be canceled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	// In main.go, before exit
	defer func() {
		if err := config.CloseDB(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}

		if err := config.CloseCache(); err != nil {
			log.Printf("Error closing cache connection: %v", err)
		}
	}()

	// HTTP Server
	api.Run(ctx)
}
