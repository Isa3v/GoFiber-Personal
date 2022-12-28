// Пакет позволяет работать с партнерским REST битрикса.
// Все методы можно посмотреть тут: https://dev.1c-bitrix.ru/learning/course/index.php?COURSE_ID=133&INDEX=Y
package marketplace

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// Базовый URL
const BaseUrlBitrixPartners string = "https://partners.1c-bitrix.ru/rest/"

type Config struct {
	// Идентификатор партнера битрикс
	PartnerId string `json:"partnerId"`

	// Код доступа партнера
	ParnerCode string `json:"partnerCode"`

	// Базовый URL
	BaseUrl string
}

type Marketplace struct {
	// Конфиг
	config Config
}

// Структура ответа Rest битрикса
type ResponseMap struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

// Вызываем маркетплейс
func New(config ...Config) (*Marketplace, error) {
	// Create a new app
	marketplace := &Marketplace{
		// Create config
		config: Config{},
	}

	// Override config if provided
	if len(config) > 0 {
		marketplace.config = config[0]
	}

	// Проверяем значения
	if len(marketplace.config.BaseUrl) == 0 {
		marketplace.config.BaseUrl = BaseUrlBitrixPartners
	}

	// Обязательные поля
	if len(marketplace.config.ParnerCode) == 0 {
		return nil, errors.New("ParnerCode is required in marketplace client")
	}

	// Обязательные поля
	if len(marketplace.config.PartnerId) == 0 {
		return nil, errors.New("ParnerCode is required in marketplace client")
	}

	// Return app
	return marketplace, nil
}

func (marketplace *Marketplace) GetConfig() Config {
	return marketplace.config
}

// Get запрос к методам
func (marketplace *Marketplace) Get(method string, params map[string]string) ResponseMap {
	// Дефолтный url
	urlParse, err := url.Parse(marketplace.config.BaseUrl)
	if err != nil {
		panic(err)
	}

	// Собираем параметры query
	queryParams := urlParse.Query()
	queryParams.Set("action", method)
	queryParams.Set("partnerId", marketplace.config.PartnerId)
	queryParams.Set("auth", marketplace.makeMD5Token(method))

	// Добавляем параметры из вызова метода
	for key, value := range params {
		queryParams.Set(key, value)
	}

	// Применяем параметры к url
	urlParse.RawQuery = queryParams.Encode()

	// Объявленем реквест и респонс
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	// Параметры запроса
	req.Header.SetConnectionClose()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentType("application/json")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", "Golang-Marketplace-Client/1.1")
	req.SetRequestURI(urlParse.String())

	// Вызываем клиент
	client := &fasthttp.Client{}
	if err := client.Do(req, res); err != nil {
		panic(err.Error())
	}

	// Получем ответ body и конвертируем из json в map
	bodyBytes := res.Body()
	var jsonResult ResponseMap
	json.Unmarshal(bodyBytes, &jsonResult)
	return jsonResult
}

// Функция собирает "токен" для запросов в маркетплейс
func (marketplace *Marketplace) makeMD5Token(method string) string {
	stringToHash := method + "|" + marketplace.config.PartnerId + "|" + marketplace.config.ParnerCode
	binHash := md5.Sum([]byte(stringToHash))
	return hex.EncodeToString(binHash[:])
}
