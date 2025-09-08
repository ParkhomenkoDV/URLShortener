package model

type Request struct {
	URL string `json:"url" validate:"required,url"`
}
