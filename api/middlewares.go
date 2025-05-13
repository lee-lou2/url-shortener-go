package api

import (
	"strings"
	"url-shortener-go/pkg"

	"github.com/gofiber/fiber/v3"
)

// jwtAuth JWT Authentication Middleware
// Validates the Authorization value in the request header to verify the JWT token's validity.
// If the token is valid, it stores the user information extracted from the token in the context.
//
// Authentication header format: Bearer <token>
//
// Process:
// 1. Check for the existence of the Authorization header and Bearer schema
// 2. Parse and validate the JWT token
// 3. Store the user information extracted from the token in the context
//
// Error responses:
//   - 401 Unauthorized:
//   - When Authorization header is missing or not using Bearer schema
//   - When JWT token is invalid
func jwtAuth(c fiber.Ctx) error {
	authHeader := c.Get(fiber.HeaderAuthorization)
	const BearerSchema = "Bearer "
	var token string
	if authHeader != "" && strings.HasPrefix(authHeader, BearerSchema) {
		token = authHeader[len(BearerSchema):]
	} else {
		token = c.Cookies("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No token provided",
			})
		}
	}
	user, err := pkg.ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: " + err.Error(),
		})
	}
	c.Locals("user_id", user)
	return c.Next()
}
