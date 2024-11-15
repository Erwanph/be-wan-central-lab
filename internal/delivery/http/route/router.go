package route

import (
	"github.com/Erwanph/be-wan-central-lab/internal/delivery/http"

	"github.com/Erwanph/be-wan-central-lab/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	AuthMiddleware *middleware.AuthMiddleware
}

func (c *RouteConfig) Setup() {
	c.App.Use(cors.New(cors.Config{
		AllowHeaders:     "Authorization,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	api := c.App.Group("api")
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).SendStatus(fiber.StatusOK)
	})

	c.SetupV1Route(api)
}

func (c *RouteConfig) SetupV1Route(api fiber.Router) {
	api = api.Group("v1")
	c.SetupAuthRoute(api)
	c.SetupProfileRoute(api)
}

func (c *RouteConfig) SetupAuthRoute(api fiber.Router) {
	auth := api.Group("auth")
	auth.Post("/register", c.UserController.RegisterUser)
	auth.Post("/verify-email", c.UserController.VerifyEmailRegister)
	auth.Post("/login", c.UserController.Login)

}
func (c *RouteConfig) SetupProfileRoute(api fiber.Router) {
	profiles := api.Group("profile")
	profiles.Use(c.AuthMiddleware.CheckSession)
	profiles.Get("/", c.UserController.GetProfiles)
	profiles.Patch("/", c.UserController.UpdateProfiles)
	profiles.Patch("/score", c.UserController.UpdateScore)
}
