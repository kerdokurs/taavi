package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"kerdo.dev/taavi/handlers"
	"kerdo.dev/taavi/middleware"
)

func GetApp() *fiber.App {
	engine := html.New("templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	app.Use(middleware.Htmx())
	app.Use(logger.New())

	handlers.Init(app)

	return app
}
