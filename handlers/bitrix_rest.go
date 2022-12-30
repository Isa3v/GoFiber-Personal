package handlers

import (
	"os"

	"isaev.digital.api/pkg/bitrix_marketplace"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

func (api *HandlerApi) GetBitrixRestCurrentParnerModules(c *fiber.Ctx) error {
	// Ключ для кеша
	cacheKey := "api:bitrixrest:marketplace.product.list:currentPartner"

	// Если нашли кеш
	if value, found := api.cache.Get(cacheKey); found {
		return c.JSON(value)
	}

	marketplaceClient, err := bitrix_marketplace.New(bitrix_marketplace.Config{
		PartnerId:  os.Getenv("BITRIX_PARTNER_ID"),
		ParnerCode: os.Getenv("BITRIX_PARTNER_CODE"),
	})

	if err != nil {
		panic(err)
	}

	// Параметры запроса
	params := map[string]string{}
	params["filter[modulePartnerId]"] = marketplaceClient.GetConfig().PartnerId

	result, err := marketplaceClient.Get("marketplace.product.list", params)

	if err != nil {
		panic(err)
	}

	// Записываем кеш
	api.cache.Set(cacheKey, result, cache.DefaultExpiration)

	// TODO при кеше что делаем?
	if len(result.Error) > 0 && len(result.Error[0].CODE) > 0 {
		c.Status(400)
		if result.Error[0].CODE == "ACTION_NOT_EXISTS" {
			c.Status(404)
		}
	}

	return c.JSON(result)
}
