package router

import (
	"github.com/gofiber/fiber/v2"
	"isaev.digital.api/handlers"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	v1 := app.Group("/api/v1")

	// Инициализуем инстанс контроллера
	apiController := handlers.New()

	// Get marketplace modules
	v1.Get("/modules/partner", apiController.GetBitrixRestCurrentParnerModules)
}
