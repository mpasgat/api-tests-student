package models

// ErrorResponse тело ответа с ошибкой, общее для всех ручек
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details any    `json:"details"`
}
