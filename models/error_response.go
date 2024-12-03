// models/error_response.go

package models

// ErrorResponse представляет структуру ответа при ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}
