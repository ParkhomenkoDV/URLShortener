package model

// Request представляет структуру входящего запроса на создание короткой ссылки.
// Используется в endpoint'е POST /api/shorten для получения оригинального URL от клиента.
//
// Пример JSON:
//
//	{"url": "https://example.com/very/long/path"}
type Request struct {
	URL string `json:"url" validate:"required,url"` // Оригинальный URL для сокращения. Должен быть валидным URL и не пустым.
}

// Response представляет структуру ответа с созданной короткой ссылкой.
// Возвращается клиенту после успешного создания короткого URL.
//
// Пример JSON:
//
//	{"result": "http://localhost:8080/abc123"}
type Response struct {
	Result string `json:"result"` // Полный короткий URL, включая базовый адрес сервиса
}

// URLRecord представляет структуру для хранения информации о ссылке в хранилище.
// Содержит полные метаданные URL, включая идентификаторы для связи с пользователем.
//
// Используется для:
//   - Сериализации/десериализации в JSON файлах
//   - Хранения в базе данных
//   - Внутренней обработки в сервисном слое
type URLRecord struct {
	ID          int    `json:"id"`           // Уникальный идентификатор записи в системе (автоинкремент)
	ShortURL    string `json:"short_url"`    // Короткий идентификатор URL (уникальный ключ)
	OriginalURL string `json:"original_url"` // Полный оригинальный URL
	UserID      string `json:"user_id"`      // Идентификатор пользователя-владельца ссылки
}

// BatchRequest представляет элемент массива запроса для пакетного создания коротких ссылок.
// Используется в endpoint'е POST /api/shorten/batch для массового сокращения URL.
//
// Поля:
//   - CorrelationID: идентификатор для сопоставления запроса и ответа
//   - OriginalURL: оригинальный URL для сокращения
//
// Пример JSON массива:
//
//	[
//	  {"correlation_id": "1", "original_url": "https://example1.com"},
//	  {"correlation_id": "2", "original_url": "https://example2.com"}
//	]
type BatchRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`   // Уникальный идентификатор для сопоставления запроса и ответа в пакете
	OriginalURL   string `json:"original_url" validate:"required,url"` // Оригинальный URL для сокращения
}

// BatchResponse представляет элемент массива ответа на пакетное создание коротких ссылок.
// Содержит результат обработки одного URL из пакетного запроса.
//
// Поля:
//   - CorrelationID: идентификатор из соответствующего запроса
//   - ShortURL: созданный короткий URL
//
// Пример JSON массива:
//
//	[
//	  {"correlation_id": "1", "short_url": "http://localhost:8080/abc123"},
//	  {"correlation_id": "2", "short_url": "http://localhost:8080/def456"}
//	]
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"` // Идентификатор из исходного BatchRequest для сопоставления
	ShortURL      string `json:"short_url"`      // Полный короткий URL для соответствующего оригинального URL
}

// UserURL представляет структуру URL пользователя для ответа в API.
// Используется в endpoint'е GET /api/user/urls для возврата списка ссылок пользователя.
//
// Поля:
//   - ShortURL: полный короткий URL (включая базовый адрес)
//   - OriginalURL: оригинальный полный URL
//
// Пример JSON:
//
//	{
//	  "short_url": "http://localhost:8080/abc123",
//	  "original_url": "https://example.com/very/long/path"
//	}
type UserURL struct {
	ShortURL    string `json:"short_url"`    // Полный короткий URL пользователя
	OriginalURL string `json:"original_url"` // Оригинальный URL, соответствующий короткой ссылке
}

// DeleteURLsRequest представляет запрос на удаление URL пользователя.
// Является псевдонимом для массива строк - списка коротких идентификаторов URL для удаления.
// Используется в endpoint'е DELETE /api/user/urls для мягкого удаления ссылок.
//
// Пример JSON:
//
//	["abc123", "def456", "ghi789"]
type DeleteURLsRequest []string
