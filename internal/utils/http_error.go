// http_error.go
package utils

type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewHTTPError(status int, message string) *HTTPError {
	return &HTTPError{
		Status:  status,
		Message: message,
	}
}
