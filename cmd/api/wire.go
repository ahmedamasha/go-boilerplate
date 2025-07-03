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
	UserController *controllers.UserController
}

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		config.NewDatabaseConnection,
		repositories.NewUserRepository,
		services.NewUserService,
		controllers.NewUserController,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}
