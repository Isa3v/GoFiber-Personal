package api

import (
	modules "isaev.digital.api/app/controller"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	v1 := app.Group("/api/v1")

	// Get marketplace modules
	v1.Get("/modules/partner", modules.GetCurrentPartnerModules)
}
