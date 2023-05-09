package router

import (
	"github.com/gofiber/fiber/v2"
	"isaev.digital.api/handlers"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	v1 := app.Group("/api/v1")

	// Инициализиуруем контроллер для работы с career.habr
	habrApi := handlers.HabrInit()
	v1.Get("/profile", habrApi.Get)

	// Инициализиуруем контроллер для работы с bitrix marketplace
	bitrixApi := handlers.BitrixInit()
	v1.Get("/modules/partner", bitrixApi.Get)
}
