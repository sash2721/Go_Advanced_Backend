package main

import (
	"advancedBackend/handlers"
	"advancedBackend/middlewares"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Apply RequestMiddleware to every route on this router
	r.Use(middlewares.RequestMiddleware)

	// Apply LoggingMiddleware to every route on this router
	r.Use(middlewares.LoggingMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "Chi Router Started"}`))
	})

	// Adding the Health route
	r.Get("/health", handlers.HealthHandler)

	// Adding the Time route
	r.Get("/time", handlers.TimeHandler)

	// Adding the Echo route
	r.Post("/echo", handlers.EchoHandler)

	// defining the PORT
	var PORT string = ":3000"

	server := &http.Server{
		Addr:         PORT,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Creating channel to listen for OS signals
	// This will create a context that will be cancelled when SIGINT or SIGTERM is recieved
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // cancel the context at the end

	// starting the server in a goroutine so it doesn't block
	go func() {
		fmt.Printf("Server listening on the PORT%s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error(
				"Error while starting the server:",
				slog.Any("Error", err),
			)
		}
	}()

	// Block here and wait for the signal
	<-ctx.Done()

	slog.Info("Shutdown Signal received, shutting down gracefully!")

	// Create a context with 5 second timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown using the shutdown context (Attempting graceful shutdown)
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error(
			"Server forced to shutdown:",
			slog.Any("Error", err),
		)
	}

	slog.Info("Server Exited!")
}
