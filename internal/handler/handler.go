package handler

import (
	"encoding/json"
	"errors"
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

// Handler — структура хендлера
type Handler struct {
	Service       *service.Service
	Configuration *config.Config
	AuthService   *auth.AuthService
}

// Конструктор хендлера
func New(
	ginEngine *gin.Engine,
	service *service.Service,
	configuration *config.Config,
) {
	// Создаем сервис аутентификации
	authService := auth.New(configuration.AuthSecretKey)
	handler := &Handler{
		Service:       service,
		Configuration: configuration,
		AuthService:   authService,
	}

	// Добавляем middleware перед регистрацией маршрутов
	ginEngine.Use(middleware.GzipMiddleware())
	ginEngine.Use(middleware.LoggingMiddleware())
	ginEngine.Use(middleware.AuthMiddleware(authService))

	// Регистрируем маршруты
	ginEngine.POST("/api/shorten", handler.SendJSONURL)
	ginEngine.POST("/api/shorten/batch", handler.SendJSONURLBatch)
	ginEngine.POST("/", handler.SendURL)
	ginEngine.GET("/:id", handler.GetURL)
	ginEngine.GET("/ping", handler.Ping)
	ginEngine.GET("/api/user/urls", handler.GetUserURLs)
}

// handleServiceError обрабатывает ошибки сервиса и отправляет соответствующий текстовый ответ
func (h *Handler) handleServiceError(c *gin.Context, err error, shortURL string) {
	if errors.Is(err, repository.ErrRowExists) {
		c.String(http.StatusConflict, h.Configuration.ShortAddress+"/"+shortURL)
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.Abort()
}

// handleServiceErrorJSON обрабатывает ошибки сервиса и отправляет соответствующий JSON ответ
func (h *Handler) handleServiceErrorJSON(c *gin.Context, err error, shortURL string) {
	if errors.Is(err, repository.ErrRowExists) {
		var response model.Response
		response.Result = h.Configuration.ShortAddress + "/" + shortURL
		c.JSON(http.StatusConflict, response)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Abort()
}

// handleGenericErrorJSON обрабатывает общие ошибки и отправляет JSON ответ
func (h *Handler) handleGenericErrorJSON(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

// handleGenericErrorText обрабатывает общие ошибки и отправляет текстовый ответ
func (h *Handler) handleGenericErrorText(c *gin.Context, statusCode int, message string) {
	c.String(statusCode, message)
	c.Abort()
}

// Обработка POST запроса: сокращение URL
func (h *Handler) SendURL(c *gin.Context) {
	// Проверка Content-Type
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
	// Получаем userID из контекста
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание короткой ссылки
	shortURL, err := h.Service.CreateShortURL(string(body), userIDStr)
	if err != nil {
		h.handleServiceError(c, err, shortURL)
		return
	}

	// Response: текст с полным URL
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, h.Configuration.ShortAddress+"/"+shortURL)
}

// Обработка POST запроса: сокращение URL (JSON)
func (h *Handler) SendJSONURL(c *gin.Context) {

	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorText(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Получаем данные из body
	var request model.Request

	// Странный тест на обязательный юз encoding/json, чисто ради галочки...
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Валидация реквеста
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем userID из контекста
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание короткой ссылки
	shortURL, err := h.Service.CreateShortURL(request.URL, userIDStr)
	if err != nil {
		h.handleServiceErrorJSON(c, err, shortURL)
		return
	}

	// Response: JSON с полным URL
	var response model.Response
	response.Result = h.Configuration.ShortAddress + "/" + shortURL

	// Пишем ответ
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
}

// Обработка POST запроса: пакетное сокращение URL (JSON)
func (h *Handler) SendJSONURLBatch(c *gin.Context) {

	// Проверка Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		h.handleGenericErrorJSON(c, http.StatusBadRequest, "Invalid ContentType, application/json only")
		return
	}

	// Получаем данные из body
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

	// Валидация всех реквестов
	validate := validator.New()
	for _, request := range requests {
		if err := validate.Struct(request); err != nil {
			h.handleGenericErrorJSON(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Извлекаем URL из реквестов
	urls := make([]string, len(requests))
	correlationMap := make(map[string]string) // originalURL -> correlationID
	for i, request := range requests {
		urls[i] = request.OriginalURL
		correlationMap[request.OriginalURL] = request.CorrelationID
	}

	// Получаем userID из контекста
	userID, _ := c.Get(middleware.UserIDKey)
	userIDStr := userID.(string)

	// Создание коротких ссылок пакетом
	shortURLsMap, err := h.Service.CreateShortURLsBatch(urls, userIDStr)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, "Error creating short URL")
		return
	}

	// Формируем ответ
	responses := make([]model.BatchResponse, 0, len(requests))
	for originalURL, shortURL := range shortURLsMap {
		if correlationID, exists := correlationMap[originalURL]; exists {
			responses = append(responses, model.BatchResponse{
				CorrelationID: correlationID,
				ShortURL:      h.Configuration.ShortAddress + "/" + shortURL,
			})
		}
	}

	// Пишем ответ
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, responses)
}

// Обработка GET запроса: редирект по короткой ссылке
func (h *Handler) GetURL(c *gin.Context) {
	// Получаем параметр из URL: /:id
	shortURL := c.Param("id")

	// Ищем полную ссылку
	fullURL, err := h.Service.GetFullURL(shortURL)
	if err != nil {
		h.handleGenericErrorText(c, http.StatusBadRequest, "URL not found")
		return
	}

	// Редирект (307)
	c.Redirect(http.StatusTemporaryRedirect, fullURL)
}

// Ping PostgreSQL
func (h *Handler) Ping(c *gin.Context) {
	if err := h.Service.PingPostgreSQL(); err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// GetUserURLs возвращает все URL пользователя
func (h *Handler) GetUserURLs(c *gin.Context) {
	// Получаем userID из контекста
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

	// Получаем URL пользователя
	userURLs, err := h.Service.GetUserURLs(userIDStr)
	if err != nil {
		h.handleGenericErrorJSON(c, http.StatusInternalServerError, "Error retrieving user URLs")
		return
	}

	// Если у пользователя нет URL
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
