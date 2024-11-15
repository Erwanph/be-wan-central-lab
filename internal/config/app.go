package config

import (
	"github.com/Erwanph/be-wan-central-lab/internal/delivery/http"
	"github.com/Erwanph/be-wan-central-lab/internal/delivery/http/middleware"
	"github.com/Erwanph/be-wan-central-lab/internal/delivery/http/route"
	"github.com/Erwanph/be-wan-central-lab/internal/repository"
	"github.com/Erwanph/be-wan-central-lab/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type BootstrapConfig struct {
	MongoDB1 *mongo.Client
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	userRepository := repository.NewUserRepository(config.MongoDB1)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.Log, config.Validate, userRepository, config.Config)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log, config.Config)

	// setup middleware
	authMiddleware := middleware.NewAuthMiddleware(config.Log, userUseCase, config.Config)
	// config.App.Use(authMiddleware.Handle)
	routeConfig := route.RouteConfig{
		App:                  config.App,
		UserController:       userController,
		AuthMiddleware:       authMiddleware,
	}
	routeConfig.Setup()
}
