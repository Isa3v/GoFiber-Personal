package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type HabrApi struct {
	fileStoragePath string     // Инициализированный кеш,
	config          HabrConfig // Конфиг подключения
}

type HabrConfig struct {
	ProfileId string `json:"profileId"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

func HabrInit() *HabrApi {
	api := &HabrApi{
		fileStoragePath: "storage/habr",
		config: HabrConfig{
			ProfileId: os.Getenv("HABR_PRODILE_ID"),
			Login:     os.Getenv("HABR_PRODILE_LOGIN"),
			Password:  os.Getenv("HABR_PRODILE_PASSWORD"),
		},
	}

	// Создаем директорию, если ее еще нет
	os.MkdirAll(api.fileStoragePath, os.ModePerm)

	// Читаем данные из файла
	_, err := ioutil.ReadFile(api.fileStoragePath + "/profile.json")
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

func (api *HabrApi) Get(c *fiber.Ctx) error {
	// Читаем данные из файла
	data, err := ioutil.ReadFile(api.fileStoragePath + "/profile.json")
	if err != nil {
		data, err = api.getData()
		if err != nil {
			panic(err)
		}
	}

	// Разбираем в JSON
	var value map[string]interface{}
	err = json.Unmarshal(data, &value)
	if err != nil {
		panic(err)
	}

	return c.JSON(value)
}

func (api *HabrApi) getData() ([]byte, error) {
	// Получаем данные по авторизации на habr
	profileData, err := api.fillData()
	if err != nil {
		return nil, err
	}

	// Создаем директорию, если ее еще нет
	os.MkdirAll(api.fileStoragePath, os.ModePerm)

	// Записываем JSON в файл
	err = ioutil.WriteFile(api.fileStoragePath+"/profile.json", profileData, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return profileData, nil
}

func (api *HabrApi) fillData() ([]byte, error) {
	// Шаг 1: Авторизация
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Устанавливаем ссылку и метод запроса
	req.SetRequestURI("https://account.habr.com/ajax/login/")
	req.Header.SetMethod("POST")

	// Устанавливаем заголовки
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://account.habr.com")
	req.Header.Set("Referer", "https://account.habr.com/login/?consumer=career&state=bslogin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	// Устаналиваем аргументы
	args.Add("state", "bslogin")
	args.Add("consumer", "career")
	args.Add("email", api.config.Login)
	args.Add("password", api.config.Password)
	args.Add("captcha", "")
	args.Add("g-recaptcha-response", "")
	args.Add("captcha_type", "recaptcha")

	req.SetBody(args.QueryString())

	log.Print("Отправили запрос на авторизацию...")

	// Делаем запрос
	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}

	log.Print("Получили успешный ответ (200)...")

	// Из запроса получаем ссылку
	body := resp.Body()
	regex := regexp.MustCompile(`window.location.href\s*=\s*'(.+)'`)
	matches := regex.FindSubmatch(body)
	if len(matches) != 2 {
		log.Printf("Status code: %d", resp.StatusCode())
		log.Printf("Response body: %s", resp.Body())
		return nil, errors.New("Не смогли получить ссылку из ответа авторизации!")
	}

	careerLinkAuth := string(matches[1])

	log.Print("Успешно спарсили ссылку для авторизации из ответа...")

	// Шаг 2: Переходим по ссылке авторизации
	req.CopyTo(req) // Копируем предыдущий запрос

	// Устанавливаем ссылку и метод запроса
	req.SetRequestURI(careerLinkAuth)
	req.Header.SetMethod("GET")

	// Делаем запрос
	if err := fasthttp.Do(req, resp); err != nil {
		return nil, errors.New("Не смогли перейти по ссылке авторизации!")
	}

	// Запрос должен вернуть 302 редирект
	if resp.StatusCode() != 302 {
		log.Printf("Status code: %d", resp.StatusCode())
		log.Printf("Response body: %s", resp.Body())
		return nil, errors.New("Ссылка авторизации не вернула редирект!")
	}

	log.Print("Успешно перешли по ссылке для авторизации...")
	cookie := resp.Header.PeekCookie("_career_session")
	if cookie == nil {
		return nil, errors.New("В ответе не найдены куки '_career_session'!")
	}

	// Шаг 3: Переходим в профиль и извлекаем данные
	req.CopyTo(req) // Копируем предыдущий запрос
	req.SetRequestURI("https://career.habr.com/" + api.config.ProfileId)
	req.Header.SetMethod("GET")
	req.Header.Add("Cookie", string(cookie))

	// Send request and get response
	if err := fasthttp.Do(req, resp); err != nil {
		log.Printf("Status code: %d", resp.StatusCode())
		log.Printf("Response body: %s", resp.Body())
		return nil, errors.New("Ошибка запроса профиля!")
	}

	log.Print("Успешно получили профиль...")

	// Поиск блока <script> с JSON
	re := regexp.MustCompile(`<script\s+type="application/json"\s+data-ssr-state="true">\s*(.+?)\s*</script>`)
	match := re.FindStringSubmatch(string(resp.Body()))
	if len(match) != 2 {
		log.Printf("Status code: %d", resp.StatusCode())
		log.Printf("Response body: %s", resp.Body())
		return nil, errors.New("Не смогли обработать JSON!")
	}
	jsonBytes := []byte(match[1])

	return jsonBytes, nil
}
