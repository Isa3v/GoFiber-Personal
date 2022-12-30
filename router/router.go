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
	apiHandlers := handlers.New()
	// Get career habr profile
	v1.Get("/profile", apiHandlers.GetHabrRestProfile)
	// Get marketplace modules
	v1.Get("/modules/partner", apiHandlers.GetBitrixRestCurrentParnerModules)
}
