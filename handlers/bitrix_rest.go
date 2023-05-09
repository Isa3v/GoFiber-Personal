package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"isaev.digital.api/pkg/bbcode"
	"isaev.digital.api/pkg/bitrix_marketplace"
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

type BitrixApi struct {
	fileStoragePath   string // Инициализированный кеш
	marketplaceClient *bitrix_marketplace.Marketplace
}

func BitrixInit() *BitrixApi {
	// Создаем подключение к Bitrix
	marketplaceClient, err := bitrix_marketplace.New(bitrix_marketplace.Config{
		PartnerId:  os.Getenv("BITRIX_PARTNER_ID"),
		ParnerCode: os.Getenv("BITRIX_PARTNER_CODE"),
	})
	if err != nil {
		panic(err)
	}

	// Инициализуем кеш
	api := &BitrixApi{
		fileStoragePath:   "storage/bitrix",
		marketplaceClient: marketplaceClient,
	}

	// Создаем директорию, если ее еще нет
	os.MkdirAll(api.fileStoragePath, os.ModePerm)

	// Читаем данные из файла
	_, err = ioutil.ReadFile(api.fileStoragePath + "/modules.json")
	if err != nil {
		_, err := api.getData()
		if err != nil {
			panic(err)
		}
	}

	// Через опредленный период обновляем кеш в файле
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			log.Print("Обновляем кеш...")
			// Получаем данные по авторизации на habr
			_, err := api.getData()
			if err != nil {
				panic(err)
			}
		}
	}()

	return api
}

func (api *BitrixApi) Get(c *fiber.Ctx) error {
	// Читаем данные из файла
	data, err := ioutil.ReadFile(api.fileStoragePath + "/modules.json")
	if err != nil {
		data, err = api.getData()
		if err != nil {
			panic(err)
		}
	}

	// Разбираем в JSON
	var value []interface{}
	err = json.Unmarshal(data, &value)
	if err != nil {
		panic(err)
	}

	return c.JSON(value)
}

func (api *BitrixApi) getData() ([]byte, error) {
	// Параметры запроса
	params := map[string]string{}
	params["filter[modulePartnerId]"] = api.marketplaceClient.GetConfig().PartnerId

	// Получаем список модулей по фильтру (params)
	var result interface{}
	resultResponse, err := api.marketplaceClient.Get("marketplace.product.list", params)
	if err != nil {
		panic(err)
	}

	// Проверяем ошибки
	if len(resultResponse.Error) > 0 && len(resultResponse.Error[0].CODE) > 0 {
		// error
	}
	// Записываем кеш, если нет ошибок
	// Форматируем результат полученный с API
	result = api.formatParnerModules(resultResponse.Result["list"].([]interface{}))

	// Разбираем в JSON
	// convert the map to JSON
	moduleData, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	// Создаем директорию, если ее еще нет
	os.MkdirAll(api.fileStoragePath, os.ModePerm)

	// Записываем JSON в файл
	err = ioutil.WriteFile(api.fileStoragePath+"/modules.json", moduleData, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return moduleData, nil
}

func (api *BitrixApi) formatParnerModules(listModules []interface{}) []BitrixPartnerModule {

	// Создаем слайс с определенным кол-вом элементов
	modulesPartnerList := make([]BitrixPartnerModule, len(listModules))

	// Компилятор BB-кодов
	bbcodeCompiler := bbcode.New()

	// Переводим в структуру каждый элемент списка модулей
	for i, val := range listModules {
		if val != nil {
			fields := val.(map[string]interface{})

			modulesPartnerList[i] = BitrixPartnerModule{
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
