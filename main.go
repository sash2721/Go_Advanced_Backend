package main

import (
	"advancedBackend/configs"
	"advancedBackend/handlers"
	"advancedBackend/middlewares"
	"advancedBackend/services"
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

	// Applying middlewares to all the routes
	r.Use(middlewares.RequestMiddleware)
	r.Use(middlewares.LoggingMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "Chi Router Started"}`))
	})

	// Creating the services
	healthService := services.NewHealthService()
	timeService := services.NewTimeService()
	echoService := services.NewEchoService()
	authService := services.NewAuthService()

	// Creating handlers and injecting the created services into them
	healthHandler := &handlers.HealthHandler{Service: healthService}
	timeHandler := &handlers.TimeHandler{Service: timeService}
	echoHandler := &handlers.EchoHandler{Service: echoService}
	authHandler := &handlers.AuthHandler{Service: authService}

	// Registering the routes
	// Public Routes
	r.Post("/auth/login", authHandler.HandleLogin)
	r.Post("/auth/signup", authHandler.HandleSignUp)
	r.Get("/health", healthHandler.HandleHealthFunction)

	// Protected Routes
	r.Group(func(r chi.Router) {
		// using the auth middleware only for protected routes
		r.Use(middlewares.AuthMiddleware)

		r.Get("/time", timeHandler.HandleTimeFunction)
		r.Post("/echo", echoHandler.HandleEchoFunction)
	})

	// getting the configs
	serverConfig := configs.GetServerConfig()

	// initialising the server
	var server *http.Server

	if serverConfig.Env == "development" {
		server = &http.Server{
			Addr:         serverConfig.Port,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		}
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
