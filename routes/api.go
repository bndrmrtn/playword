package routes

import (
	"github.com/bndrmrtn/playword/handlers"
	"github.com/gofiber/fiber/v2"
)

func AddApiRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/game", handlers.GameHandler)
}
