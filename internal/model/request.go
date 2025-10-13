package model

// Model Request
type Request struct {
	URL string `json:"url" validate:"required,url"`
}

// Model Response
type Response struct {
	Result string `json:"result"`
}

// Model for URL storage in JSON format
type URLRecord struct {
	ID          int    `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Model for batch request
type BatchRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"original_url" validate:"required,url"`
}

// Model for batch response
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
