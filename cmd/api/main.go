package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cusror_ai/internal/config"
	"cusror_ai/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Redis connection for side effect/logging
	_, err = config.NewRedisConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize dependencies using Wire
	app, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// Use CORS middleware from internal/middleware
	router.Use(middleware.CORS)

	// Setup routes
	setupRoutes(router, app)

	// Create server
	server := &http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", cfg.GetServerAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRoutes(router *mux.Router, app *App) {
	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	apiV1.HandleFunc("/users", app.UserController.GetAllUsers).Methods("GET")
	apiV1.HandleFunc("/users", app.UserController.CreateUser).Methods("POST")
	apiV1.HandleFunc("/users/{id}", app.UserController.GetUserByID).Methods("GET")
	apiV1.HandleFunc("/users/{id}", app.UserController.UpdateUser).Methods("PUT")
	apiV1.HandleFunc("/users/{id}", app.UserController.DeleteUser).Methods("DELETE")

	// Dashboard endpoint (protected)
	router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(app.UserController.Dashboard))).Methods("GET")

	// Login endpoint
	router.HandleFunc("/login", app.UserController.Login).Methods("POST")

	// Health check
	router.HandleFunc("/health", app.UserController.HealthCheck).Methods("GET")
}
