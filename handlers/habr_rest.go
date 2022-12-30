package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

func (api *HandlerApi) GetHabrRestProfile(c *fiber.Ctx) error {
	// Ключ для кеша
	cacheKey := "api:habr:currentProfile"

	// Если нашли кеш
	if value, found := api.cache.Get(cacheKey); found {
		return c.JSON(value)
	}

	api.cache.Set(cacheKey, "bar", cache.DefaultExpiration)

	return c.JSON("bar")
}
