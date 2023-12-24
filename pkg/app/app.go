package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"kerdo.dev/taavi/pkg/handlers"
	"kerdo.dev/taavi/pkg/middleware"
)

func GetApp() *fiber.App {
	app := fiber.New()
	app.Static("/", "./public")

	app.Use(middleware.Htmx())
	app.Use(logger.New())
	app.Use(basicauth.New(basicauth.Config{
		Users: loadAuthUsers(),
	}))

	handlers.Init(app)

	return app
}
