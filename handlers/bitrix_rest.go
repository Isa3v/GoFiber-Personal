package handlers

import (
	"os"
	"strconv"

	"isaev.digital.api/pkg/bbcode"
	"isaev.digital.api/pkg/bitrix_marketplace"

	"github.com/gofiber/fiber/v2"
)

type BitrixPartnerModule struct {
	BasePrice               float64       `json:"basePriceRub"`
	Price                   float64       `json:"priceRub"`
	Categories              []interface{} `json:"categories"`
	Code                    string        `json:"code"`
	CompatibleEditions      []interface{} `json:"compatibleEditions"`
	DatePublish             string        `json:"datePublish"`
	DemoUrl                 string        `json:"demoUrl"`
	Description             string        `json:"description"`
	Includes                []interface{} `json:"includes"`
	InstallationDescription string        `json:"installationDescription"`
	LogoUrl                 string        `json:"logoUrl"`
	ModulePartnerId         string        `json:"modulePartnerId"`
	Name                    string        `json:"name"`
	Screenshots             []interface{} `json:"screenshots"`
	SupportDescription      string        `json:"supportDescription"`
	VideoUrl                string        `json:"videoUrl"`
}

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

	resultResponse, err := marketplaceClient.Get("marketplace.product.list", params)

	if err != nil {
		panic(err)
	}

	// Форматируем результат полученный с API
	result := api.formatParnerModules(resultResponse.Result["list"].([]interface{}))

	// Записываем кеш
	api.cache.Set(cacheKey, result, 0)

	// TODO при кеше что делаем?
	if len(resultResponse.Error) > 0 && len(resultResponse.Error[0].CODE) > 0 {
		c.Status(400)
		if resultResponse.Error[0].CODE == "ACTION_NOT_EXISTS" {
			c.Status(404)
		}
	}

	return c.JSON(result)
}

func (api *HandlerApi) formatParnerModules(listModules []interface{}) []BitrixPartnerModule {

	// Создаем слайс с определенным кол-вом элементов
	modulesPartnerList := make([]BitrixPartnerModule, len(listModules))
	bbcodeCompiler := bbcode.New()

	for i, val := range listModules {
		if val != nil {
			fields := val.(map[string]interface{})

			modulesPartnerList[i] = BitrixPartnerModule{
				BasePrice:               0,
				Price:                   0,
				Categories:              fields["categories"].([]interface{}),
				Code:                    fields["code"].(string),
				CompatibleEditions:      fields["compatibleEditions"].([]interface{}),
				DatePublish:             fields["datePublish"].(string),
				DemoUrl:                 fields["demoUrl"].(string),
				Description:             bbcodeCompiler.Compile(fields["description"].(string)),
				Includes:                fields["includes"].([]interface{}),
				InstallationDescription: bbcodeCompiler.Compile(fields["installationDescription"].(string)),
				LogoUrl:                 fields["logoUrl"].(string),
				ModulePartnerId:         fields["modulePartnerId"].(string),
				Name:                    fields["name"].(string),
				Screenshots:             fields["screenshots"].([]interface{}),
				SupportDescription:      bbcodeCompiler.Compile(fields["supportDescription"].(string)),
				VideoUrl:                fields["videoUrl"].(string),
			}

			// Приводим базовую цену к float
			switch v := fields["basePriceRub"].(type) {
			case float64:
				// v is a float64 here, so e.g. v + 1.0 is possible.
				modulesPartnerList[i].BasePrice = v
			case string:
				// v is a string here, so e.g. v + " Yeah!" is possible.
				if s, err := strconv.ParseFloat(fields["basePriceRub"].(string), 32); err == nil {
					modulesPartnerList[i].BasePrice = s
				} else {
					panic(err)
				}
			}

			// Приводим цену к float
			switch v := fields["priceRub"].(type) {
			case float64:
				// v is a float64 here, so e.g. v + 1.0 is possible.
				modulesPartnerList[i].Price = v
			case string:
				// v is a string here, so e.g. v + " Yeah!" is possible.
				if s, err := strconv.ParseFloat(fields["priceRub"].(string), 32); err == nil {
					modulesPartnerList[i].Price = s
				} else {
					panic(err)
				}
			}
		}
	}

	return modulesPartnerList
}
