package api

import (
	"context"
	"fmt"
	"log"
	"url-shortener-go/config"

	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/template/html/v2"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	jsoniter "github.com/json-iterator/go"
)

// Run Function to run the API server
//
// Initializes the Fiber web framework, configures middleware, and starts the server.
// Safely shuts down the server when the context is canceled.
//
// Main features:
//   - HTML template engine setup
//   - Fiber app instance creation and configuration
//   - Middleware setup (pprof, recover, logger, encryptcookie)
//   - Route configuration
//   - Asynchronous server execution
//   - Safe shutdown when context is canceled
//
// Parameters:
//   - ctx: Context that controls the server lifecycle
//
// Usage example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	api.Run(ctx)
func Run(ctx context.Context) {
	engine := html.New("./views", ".html")
	app := fiber.New(
		fiber.Config{
			AppName: "url-shortener",
			Views:   engine,
			// JSON Encoder and Decoder(Fast JSON library)
			JSONEncoder: jsoniter.Marshal,
			JSONDecoder: jsoniter.Unmarshal,
		},
	)

	// Middleware
	app.Use(pprof.New())
	app.Use(recoverer.New())
	secretKey := config.GetEnv("ENCRYPT_COOKIE_KEY")
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: secretKey,
	}))
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${pid} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Routes
	setV1Routes(app)

	go func() {
		log.Fatal(app.Listen(fmt.Sprintf(":%s", config.GetEnv("SERVER_PORT", "3000"))))
	}()

	<-ctx.Done()
	app.Shutdown()
}
