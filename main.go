package main

import (
	"go-auth/config"
	"go-auth/database"
	"go-auth/routes"
	"go-auth/sessions"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	// load environment variables from a `.env` file
	config.MustLoadEnvVariables()

	// connect to a postgres database and start migrations
	database.MustInit()

	// initialize the redis sessions store
	sessions.MustInitRedisStore()
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FRONT_URL"),
		AllowCredentials: true,
	}))

	app.Use(logger.New())

	routes.AttachAuthRoutesV1(app)

	log.Fatalln(app.Listen(":5050"))
}
