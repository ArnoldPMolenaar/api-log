package main

import (
	"api-log/main/src/configs"
	"api-log/main/src/database"
	"api-log/main/src/middleware"
	"api-log/main/src/routes"
	"api-log/main/src/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Register Fiber's middleware for app.
	middleware.FiberMiddleware(app)

	// Open database connection.
	if err := database.OpenDBConnection(); err != nil {
		panic(fmt.Sprintf("Could not connect to the database: %v", err))
	}

	// Register a route for API Docs (Scalar).
	// TODO: routes.ScalarRoute(app)
	// Register a private routes for app.
	routes.PrivateRoutes(app)
	// Register route for 404 Error.
	routes.NotFoundRoute(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
