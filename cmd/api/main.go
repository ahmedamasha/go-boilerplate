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

	// Start WebSocket service
	go app.WebSocketService.Run()
	go app.WebSocketService.StartHeartbeat()

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

	// Event routes
	apiV1.HandleFunc("/events", app.EventController.CreateEvent).Methods("POST")
	apiV1.HandleFunc("/events/{id}", app.EventController.GetEvent).Methods("GET")
	apiV1.HandleFunc("/events/user/{user_id}", app.EventController.GetUserEvents).Methods("GET")
	apiV1.HandleFunc("/events/type/{type}", app.EventController.GetEventsByType).Methods("GET")
	apiV1.HandleFunc("/events/recent", app.EventController.GetRecentEvents).Methods("GET")
	apiV1.HandleFunc("/events/user/{user_id}/stats", app.EventController.GetUserEventStats).Methods("GET")
	apiV1.HandleFunc("/events/simulate/shopify", app.EventController.SimulateShopifyEvent).Methods("POST")

	// Segment routes
	apiV1.HandleFunc("/segments", app.SegmentController.GetAllSegments).Methods("GET")
	apiV1.HandleFunc("/segments", app.SegmentController.CreateSegment).Methods("POST")
	apiV1.HandleFunc("/segments/{id}", app.SegmentController.GetSegment).Methods("GET")
	apiV1.HandleFunc("/segments/{id}", app.SegmentController.UpdateSegment).Methods("PUT")
	apiV1.HandleFunc("/segments/{id}", app.SegmentController.DeleteSegment).Methods("DELETE")
	apiV1.HandleFunc("/segments/analyze/{user_id}", app.SegmentController.AnalyzeUserSegments).Methods("POST")
	apiV1.HandleFunc("/segments/user/{user_id}", app.SegmentController.GetUserSegments).Methods("GET")
	apiV1.HandleFunc("/segments/{id}/users", app.SegmentController.GetSegmentUsers).Methods("GET")
	apiV1.HandleFunc("/segments/user/{user_id}/interests", app.SegmentController.GetUserInterests).Methods("GET")
	apiV1.HandleFunc("/segments/predefined", app.SegmentController.CreatePredefinedSegments).Methods("POST")
	apiV1.HandleFunc("/segments/rules/template", app.SegmentController.GetSegmentRulesTemplate).Methods("GET")

	// Offer routes
	apiV1.HandleFunc("/offers", app.OfferController.GetActiveOffers).Methods("GET")
	apiV1.HandleFunc("/offers/{id}", app.OfferController.GetOffer).Methods("GET")
	apiV1.HandleFunc("/offers/user/{user_id}", app.OfferController.GetUserOffers).Methods("GET")
	apiV1.HandleFunc("/offers/{id}/use", app.OfferController.UseOffer).Methods("POST")
	apiV1.HandleFunc("/offers/stats", app.OfferController.GetOfferStats).Methods("GET")

	// WebSocket routes
	router.HandleFunc("/ws", app.WebSocketController.HandleWebSocket).Methods("GET")
	apiV1.HandleFunc("/ws/stats", app.WebSocketController.GetConnectionStats).Methods("GET")

	// Dashboard endpoint (protected)
	router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(app.UserController.Dashboard))).Methods("GET")

	// Login endpoint
	router.HandleFunc("/login", app.UserController.Login).Methods("POST")

	// Health check
	router.HandleFunc("/health", app.UserController.HealthCheck).Methods("GET")
}
