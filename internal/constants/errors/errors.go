package errors

// Коды ошибок из поля error ответа
const (
	InvalidRequest = "invalid_request"
	Unauthorized   = "unauthorized"
	NotFound       = "not_found"
	Forbidden      = "forbidden"
	AlreadyExists  = "already_exists"
)

// Тексты ошибок авторизации из поля message ответа
const (
	AuthHeaderRequiredMessage   = "Authorization header required"
	InvalidAuthorizationMessage = "Invalid or expired token"
)
