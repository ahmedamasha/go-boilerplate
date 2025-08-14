//go:build wireinject
// +build wireinject

package main

import (
	"cusror_ai/internal/config"
	"cusror_ai/internal/controllers"
	"cusror_ai/internal/repositories"
	"cusror_ai/internal/services"

	"github.com/google/wire"
)

type App struct {
	UserController      *controllers.UserController
	EventController     *controllers.EventController
	SegmentController   *controllers.SegmentController
	OfferController     *controllers.OfferController
	WebSocketController *controllers.WebSocketController
	WebSocketService    *services.WebSocketService
}

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		// Database connection
		config.NewDatabaseConnection,

		// Repositories
		repositories.NewUserRepository,
		repositories.NewEventRepository,
		repositories.NewSegmentRepository,
		repositories.NewOfferRepository,

		// Services
		services.NewUserService,
		services.NewSegmentService,
		services.NewOfferService,
		services.NewWebSocketService,
		NewEventServiceWithDeps,

		// Controllers
		controllers.NewUserController,
		controllers.NewEventController,
		controllers.NewSegmentController,
		controllers.NewOfferController,
		controllers.NewWebSocketController,

		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

// NewEventServiceWithDeps creates EventService and sets up dependencies to avoid circular imports
func NewEventServiceWithDeps(
	eventRepo *repositories.EventRepository,
	segmentRepo *repositories.SegmentRepository,
	offerRepo *repositories.OfferRepository,
	segmentService *services.SegmentService,
	offerService *services.OfferService,
	wsService *services.WebSocketService,
) *services.EventService {
	eventService := services.NewEventService(eventRepo, segmentRepo, offerRepo)
	eventService.SetServices(segmentService, offerService, wsService)
	offerService.SetWebSocketService(wsService)
	return eventService
}
 