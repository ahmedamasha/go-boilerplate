package router

import (
	"github.com/gin-gonic/gin"
	"github.com/refda/backend/internal/app/handlers"
	"github.com/refda/backend/internal/pkg/config"
	jwtpkg "github.com/refda/backend/internal/pkg/jwt"
	"github.com/refda/backend/internal/pkg/middleware"
)

type Handlers struct {
	Auth  *handlers.AuthHandler
	User  *handlers.UserHandler
	Event *handlers.EventHandler
	Gift  *handlers.GiftHandler
	Admin *handlers.AdminHandler
}

func Setup(cfg *config.Config, jwtMgr *jwtpkg.Manager, h *Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.Static("/uploads", cfg.Storage.LocalPath)
	r.GET("/health", handlers.Health)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Auth.Login)
			auth.POST("/register", h.Auth.Register)
			auth.POST("/verify", h.Auth.Verify)
		}

		protected := v1.Group("")
		protected.Use(middleware.Auth(jwtMgr))
		{
			users := protected.Group("/users")
			{
				users.GET("/me", h.User.GetProfile)
				users.PUT("/me", h.User.UpdateProfile)
			}

			protected.POST("/events", h.Event.Create)
			protected.POST("/events/:id/withdraw", h.Event.Withdraw)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AdminAuth(cfg))
		{
			admin.POST("/events/:id/review", h.Admin.ReviewEvent)
		}

		// Public reads; optional JWT when ?mine=true
		eventsRead := v1.Group("")
		eventsRead.Use(middleware.OptionalAuth(jwtMgr))
		eventsRead.GET("/events", h.Event.GetList)
		eventsRead.GET("/events/:id", h.Event.GetByID)

		gifts := v1.Group("/gifts")
		{
			gifts.POST("/contribute", h.Gift.Contribute)
			gifts.POST("/pay", h.Gift.Pay)
		}
	}

	return r
}
