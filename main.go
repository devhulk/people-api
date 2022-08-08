package main

import (
	"appraisals-api/configs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "This is initial data"})
	})

	app.Listen(":6000")
}
