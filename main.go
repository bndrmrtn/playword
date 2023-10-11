package main

import (
	"github.com/bndrmrtn/playword/database"
	"github.com/bndrmrtn/playword/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {
	app := fiber.New()

	database.ConnectDb()

	//app.Use(recover.New())
	app.Use(logger.New())
	app.Use(helmet.New())

	routes.AddApiRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
