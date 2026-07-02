package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/refda/backend/internal/app/handlers"
	"github.com/refda/backend/internal/app/repositories"
	"github.com/refda/backend/internal/app/router"
	"github.com/refda/backend/internal/app/services"
	"github.com/refda/backend/internal/pkg/config"
	"github.com/refda/backend/internal/pkg/database"
	jwtpkg "github.com/refda/backend/internal/pkg/jwt"
	"github.com/refda/backend/internal/pkg/otp"
	"github.com/refda/backend/internal/pkg/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	store, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	jwtMgr := jwtpkg.NewManager(cfg)
	otpRepo := repositories.NewOTPRepository(db)
	otpProvider := otp.NewProvider(otpRepo, cfg)

	userRepo := repositories.NewUserRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	giftRepo := repositories.NewGiftRepository(db)
	contributionRepo := repositories.NewContributionRepository(db)
	withdrawalRepo := repositories.NewWithdrawalRepository(db)

	authSvc := services.NewAuthService(userRepo, otpProvider, jwtMgr)
	userSvc := services.NewUserService(userRepo, store)
	eventSvc := services.NewEventService(eventRepo, contributionRepo, withdrawalRepo, store)
	giftSvc := services.NewGiftService(giftRepo, eventRepo, contributionRepo, eventSvc)

	h := &router.Handlers{
		Auth:  handlers.NewAuthHandler(authSvc),
		User:  handlers.NewUserHandler(userSvc),
		Event: handlers.NewEventHandler(eventSvc),
		Gift:  handlers.NewGiftHandler(giftSvc),
		Admin: handlers.NewAdminHandler(eventSvc),
	}

	engine := router.Setup(cfg, jwtMgr, h)

	srv := &http.Server{
		Addr:    cfg.ServerAddr(),
		Handler: engine,
	}

	go func() {
		log.Printf("Refda API listening on %s", cfg.ServerAddr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}
