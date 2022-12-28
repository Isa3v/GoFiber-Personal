package main

import (
	"boilerplate/marketplace"

	"flag"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork:     *prod,          // go run app.go -prod
		JSONEncoder: json.Marshal,   // go-json
		JSONDecoder: json.Unmarshal, // go-json
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	v1.Get("/modules", func(c *fiber.Ctx) error {
		marketplaceClient, err := marketplace.New(marketplace.Config{
			PartnerId:  "10338950",
			ParnerCode: "2011295",
		})

		if err != nil {
			log.Fatalf("%+v", err)
		}

		params := make(map[string]string)
		params["filter[modulePartnerId]"] = marketplaceClient.GetConfig().PartnerId

		return c.JSON(fiber.Map{
			"result": marketplaceClient.Get("marketplace.product.list", params),
		})
	})

	// Setup static files
	app.Static("/", "./static/public")

	// Handle not founds
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendFile("./static/private/404.html")
	})

	// Listen on port 3000
	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
}
