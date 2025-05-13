package api

import (
	"github.com/gofiber/fiber/v3"
)

// setV1Routes V1 Routes
func setV1Routes(app *fiber.App) {
	// Template
	app.Get("/", indexHandler)
	app.Get("/:short_key", redirectToOriginalHandler)

	// API
	v1 := app.Group("/v1")
	{
		v1.Post("/urls", createShortUrlHandler, jwtAuth)
	}
}
