package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ParkhomenkoDV/URLShortener/internal/auth"
	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/middleware"
	"github.com/ParkhomenkoDV/URLShortener/internal/model"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/ParkhomenkoDV/URLShortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// Handler — структура хендлера HTTP-запросов для сервиса сокращения URL.
// Содержит зависимости: сервисный слой, конфигурацию и сервис аутентификации.
type Handler struct {
	Service       *service.Service  // Бизнес-логика приложения
	Configuration *config.Config    // Конфигурационные параметры
	AuthService   *auth.AuthService // Сервис аутентификации и авторизации
}

// New инициализирует обработчики HTTP-запросов и регистрирует маршруты.
// Настраивает middleware и привязывает обработчики к соответствующим endpoint'ам.
//
// Параметры:
//   - ginEngine: роутер Gin
//   - service: сервисный слой с бизнес-логикой
//   - configuration: конфигурация приложения
func New(
	ginEngine *gin.Engine,
	service *service.Service,
	configuration *config.Config,
) {
	// Инициализируем сервис аутентификации с секретным ключом из конфигурации
	authService := auth.New(configuration.AuthSecretKey)
	handler := &Handler{
		Service:       service,
		Configuration: configuration,
		AuthService:   authService,
	}

	// Регистрируем middleware в порядке выполнения (сверху вниз)
	ginEngine.Use(middleware.GzipMiddleware())            // Сжатие HTTP-ответов
	ginEngine.Use(middleware.LoggingMiddleware())         // Логирование запросов
	ginEngine.Use(middleware.AuthMiddleware(authService)) // Аутентификация пользователей

	// Регистрируем маршруты API
	ginEngine.POST("/api/shorten", handler.SendJSONURL)            // Создание короткой ссылки (JSON)
	ginEngine.POST("/api/shorten/batch", handler.SendJSONURLBatch) // Пакетное создание коротких ссылок
	ginEngine.POST("/", handler.SendURL)                           // Создание короткой ссылки (текст)
	ginEngine.GET("/:id", handler.GetURL)                          // Редирект по короткой ссылке
	ginEngine.GET("/ping", handler.Ping)                           // Проверка доступности БД
	ginEngine.GET("/api/user/urls", handler.GetUserURLs)           // Получение URL пользователя
	ginEngine.DELETE("/api/user/urls", handler.DeleteUserURLs)     // Мягкое удаление URL пользователя
}

// handleServiceError обрабатывает ошибки сервисного слоя и отправляет соответствующий текстовый ответ.
// Особый случай: конфликт при создании дублирующейся ссылки (Status 409).
//
// Параметры:
//   - c: контекст Gin
//   - err: ошибка из сервиса
//   - shortURL: короткий идентификатор URL (для случая конфликта)
func (h *Handler) handleServiceError(c *gin.Context, err error, shortURL string) {
	if errors.Is(err, repository.ErrRowExists) {
		// Возвращаем существующую короткую ссылку с кодом 409 Conflict
		c.String(http.StatusConflict, h.Configuration.ShortAddress+"/"+shortURL)
	} else {
		// Внутренняя ошибка сервера
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.Abort()
}

// handleServiceErrorJSON обрабатывает ошибки сервисного слоя и отправляет соответствующий JSON ответ.
// Аналогично handleServiceError, но возвращает ответ в формате JSON.
func (h *Handler) handleServiceErrorJSON(c *gin.Context, err error, shortURL string) {
	if errors.Is(err, repository.ErrRowExists) {
		// Конфликт: ссылка уже существует
		var response model.Response
		response.Result = h.Configuration.ShortAddress + "/" + shortURL
		c.JSON(http.StatusConflict, response)
	} else {
		// Внутренняя ошибка сервера
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Abort()
}

// handleGenericErrorJSON обрабатывает общие ошибки валидации и отправляет JSON ответ.
// Используется для ошибок клиента (4xx) с понятным сообщением.
func (h *Handler) handleGenericErrorJSON(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

// handleGenericErrorText обрабатывает общие ошибки и отправляет текстовый ответ.
// Используется для endpoint'ов, работающих с text/plain.
func (h *Handler) handleGenericErrorText(c *gin.Context, statusCode int, message string) {
	c.String(statusCode, message)
	c.Abort()
}

// SendURL обрабатывает POST запрос для создания короткой ссылки из текстового тела запроса.
// Content-Type: text/plain
//
// Пример запроса:
//
//	POST / HTTP/1.1
//	Content-Type: text/plain
//	https://example.com/very/long/url
//
// Пример ответа:
//
//	HTTP/1.1 201 Created
//	Content-Type: text/plain
//	http://localhost:8080/abc123
func (h *Handler) SendURL(c *gin.Context) {
	// Проверка Content-Type: должен быть text/plain
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "text/plain") {
		h.handleGenericErrorText(c, http.StatusBadRequest, "Invalid ContentType, text/plain only")
		return
	}

	// Чтение тела запроса
	body, err := c.GetRawData()
	if err != nil {
		h.handleGenericErrorText(c, http.StatusBadRequest, "Error reading request body")
		return
	}

	// Получаем userID из контекста (устанавливается в middleware аутентификации)
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание короткой ссылки через сервисный слой
	shortURL, err := h.Service.CreateShortURL(string(body), userIDStr)
	if err != nil {
		h.handleServiceError(c, err, shortURL)
		return
	}

	// Успешный ответ: полный короткий URL
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, h.Configuration.ShortAddress+"/"+shortURL)
}

// SendJSONURL обрабатывает POST запрос для создания короткой ссылки из JSON тела запроса.
// Content-Type: application/json
//
// Пример запроса:
//
//	POST /api/shorten HTTP/1.1
//	Content-Type: application/json
//	{"url": "https://example.com/very/long/url"}
//
// Пример ответа:
//
//	HTTP/1.1 201 Created
//	Content-Type: application/json
//	{"result": "http://localhost:8080/abc123"}
func (h *Handler) SendJSONURL(c *gin.Context) {
	// Проверка Content-Type: должен быть application/json
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorText(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Декодируем JSON тело запроса
	var request model.Request
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация структуры запроса
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем userID из контекста аутентификации
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание короткой ссылки через сервисный слой
	shortURL, err := h.Service.CreateShortURL(request.URL, userIDStr)
	if err != nil {
		h.handleServiceErrorJSON(c, err, shortURL)
		return
	}

	// Формируем JSON ответ
	var response model.Response
	response.Result = h.Configuration.ShortAddress + "/" + shortURL

	// Успешный ответ
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
}

// SendJSONURLBatch обрабатывает POST запрос для пакетного создания коротких ссылок.
// Принимает массив URL и возвращает массив соответствующих коротких ссылок.
// Content-Type: application/json
//
// Пример запроса:
//
//	POST /api/shorten/batch HTTP/1.1
//	Content-Type: application/json
//	[
//	  {"correlation_id": "1", "original_url": "https://example1.com"},
//	  {"correlation_id": "2", "original_url": "https://example2.com"}
//	]
//
// Пример ответа:
//
//	HTTP/1.1 201 Created
//	Content-Type: application/json
//	[
//	  {"correlation_id": "1", "short_url": "http://localhost:8080/abc123"},
//	  {"correlation_id": "2", "short_url": "http://localhost:8080/def456"}
//	]
func (h *Handler) SendJSONURLBatch(c *gin.Context) {
	// Проверка Content-Type: должен быть application/json
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Декодируем массив запросов из JSON тела
	var requests []model.BatchRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&requests); err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем, что пакет не пустой
	if len(requests) == 0 {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Empty batch not allowed")
		return
	}

	// Валидация всех запросов в пакете
	validate := validator.New()
	for _, request := range requests {
		if err := validate.Struct(request); err != nil {
			h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Извлекаем URL из запросов и создаем маппинг URL -> correlation_id
	urls := make([]string, len(requests))
	correlationMap := make(map[string]string) // originalURL -> correlationID
	for i, request := range requests {
		urls[i] = request.OriginalURL
		correlationMap[request.OriginalURL] = request.CorrelationID
	}

	// Получаем userID из контекста аутентификации
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Пакетное создание коротких ссылок через сервисный слой
	shortURLsMap, err := h.Service.CreateShortURLsBatch(urls, userIDStr)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, "Error creating short URL")
		return
	}

	// Формируем ответы, сохраняя correlation_id
	responses := make([]model.BatchResponse, 0, len(requests))
	for originalURL, shortURL := range shortURLsMap {
		if correlationID, exists := correlationMap[originalURL]; exists {
			responses = append(responses, model.BatchResponse{
				CorrelationID: correlationID,
				ShortURL:      h.Configuration.ShortAddress + "/" + shortURL,
			})
		}
	}

	// Успешный ответ с массивом результатов
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, responses)
}

// GetURL обрабатывает GET запрос для редиректа по короткой ссылке.
// Извлекает короткий идентификатор из URL параметра и выполняет редирект на оригинальный URL.
// Возвращает 410 Gone если ссылка была удалена.
//
// Пример запроса:
//
//	GET /abc123 HTTP/1.1
//
// Пример ответа:
//
//	HTTP/1.1 307 Temporary Redirect
//	Location: https://example.com/very/long/url
func (h *Handler) GetURL(c *gin.Context) {
	// Получаем короткий идентификатор из URL параметра :id
	shortURL := c.Param("id")

	// Проверяем, не удалена ли ссылка (мягкое удаление)
	if deleted, err := h.Service.IsDeleted(shortURL); err == nil && deleted {
		// Ссылка удалена - возвращаем 410 Gone
		c.Status(http.StatusGone)
		return
	}

	// Ищем полный URL по короткому идентификатору
	fullURL, err := h.Service.GetFullURL(shortURL)
	if err != nil {
		h.handleGenericErrorText(c, http.StatusBadRequest, "URL not found")
		return
	}

	// Выполняем временный редирект (307)
	c.Redirect(http.StatusTemporaryRedirect, fullURL)
}

// Ping обрабатывает GET запрос для проверки доступности базы данных.
// Используется для health-check'ов и мониторинга.
//
// Пример запроса:
//
//	GET /ping HTTP/1.1
//
// Пример ответа:
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//	{"status": "OK"}
func (h *Handler) Ping(c *gin.Context) {
	if err := h.Service.PingPostgreSQL(); err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// GetUserURLs обрабатывает GET запрос для получения всех сокращенных URL пользователя.
// Требует аутентификации. Возвращает 204 No Content если у пользователя нет URL.
//
// Пример запроса:
//
//	GET /api/user/urls HTTP/1.1
//	Cookie: user_id=<signed-cookie>
//
// Пример ответа:
//
//	HTTP/1.1 200 OK
//	Content-Type: application/json
//	[
//	  {
//	    "short_url": "http://localhost:8080/abc123",
//	    "original_url": "https://example.com/very/long/url"
//	  }
//	]
func (h *Handler) GetUserURLs(c *gin.Context) {
	// Получаем userID из контекста (устанавливается middleware аутентификации)
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDStr := userID.(string)
	if userIDStr == "" {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Получаем URL пользователя через сервисный слой
	userURLs, err := h.Service.GetUserURLs(userIDStr)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, "Error retrieving user URLs")
		return
	}

	// Если у пользователя нет URL, возвращаем 204 No Content
	if len(userURLs) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Формируем ответ с полными URL
	response := make([]model.UserURL, len(userURLs))
	for i, urlData := range userURLs {
		response[i] = model.UserURL{
			ShortURL:    h.Configuration.ShortAddress + "/" + urlData["short_url"],
			OriginalURL: urlData["original_url"],
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

// DeleteUserURLs обрабатывает DELETE запрос для пометки URL как удаленных (мягкое удаление).
// Выполняется асинхронно - возвращает 202 Accepted сразу после принятия запроса.
// Content-Type: application/json
//
// Пример запроса:
//
//	DELETE /api/user/urls HTTP/1.1
//	Content-Type: application/json
//	Cookie: user_id=<signed-cookie>
//	["abc123", "def456"]
//
// Пример ответа:
//
//	HTTP/1.1 202 Accepted
func (h *Handler) DeleteUserURLs(c *gin.Context) {
	// Проверка Content-Type: должен быть application/json
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Получаем userID из контекста аутентификации
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDStr := userID.(string)
	if userIDStr == "" {
		h.handleGenericErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Декодируем JSON тело запроса (массив коротких идентификаторов)
	var deleteRequest model.DeleteURLsRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&deleteRequest); err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем, что список не пустой
	if len(deleteRequest) == 0 {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Empty list not allowed")
		return
	}

	// Запускаем асинхронное удаление URL (не блокируем клиента)
	go func() {
		if err := h.Service.DeleteURLsBatch(deleteRequest, userIDStr); err != nil {
			log.Printf("Error deleting URLs batch for user %s: %v", userIDStr, err)
		}
	}()

	// Возвращаем статус 202 Accepted - запрос принят в обработку
	c.Status(http.StatusAccepted)
}
