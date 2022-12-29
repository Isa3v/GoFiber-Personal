package main

import (
	"isaev.digital.api/router"

	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// .env config
	err := godotenv.Load()
	if err != nil {
		panic(`".env" file not found`)
	}

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	// Конфигурация recover middleware (обработка panic)
	ConfigRecover := recover.Config{
		Next:              nil,
		EnableStackTrace:  true, // enable console log trace panic
		StackTraceHandler: recover.ConfigDefault.StackTraceHandler,
	}
	app.Use(recover.New(ConfigRecover))
	app.Use(logger.New())

	// Router
	// Старт роутера
	router.SetupRoutes(app)

	// Обработка 404
	app.Use(func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return c.JSON(map[string]string{
			"message": "Not Found",
		})
	})

	// Чтение порта
	log.Fatal(app.Listen(*port))
}
