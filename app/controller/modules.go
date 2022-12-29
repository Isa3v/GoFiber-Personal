package controller

import (
	"os"

	marketplace "isaev.digital.api/app/service"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentPartnerModules(c *fiber.Ctx) error {
	marketplaceClient := marketplace.New(marketplace.Config{
		PartnerId:  os.Getenv("BITRIX_PARTNER_ID"),
		ParnerCode: os.Getenv("BITRIX_PARTNER_CODE"),
	})

	params := map[string]string{}
	params["filter[modulePartnerId]"] = marketplaceClient.GetConfig().PartnerId

	result := marketplaceClient.Get("marketplace.product.list", params)

	// TODO: завернуть в функцию
	if len(result.Error) > 0 && len(result.Error[0].CODE) > 0 {
		if result.Error[0].CODE == "ACTION_NOT_EXISTS" {
			c.Status(404)
		} else {
			c.Status(400)
		}

		return c.JSON(map[string]interface{}{
			"message": result.Error[0].CODE,
		})
	}

	return c.JSON(result)
}
