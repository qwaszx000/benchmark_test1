package main

import (
	"github.com/gofiber/fiber/v2"
)

func test_handler(ctx *fiber.Ctx) error {
	ctx.SendStatus(200)
	ctx.SendString("Hello world!")
	return nil
}

func main() {
	app := fiber.New()

	app.Get("/test_plain", test_handler)

	app.Listen("127.0.0.1:8080")
}
