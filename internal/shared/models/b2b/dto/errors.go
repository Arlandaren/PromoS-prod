package dto

import "errors"

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrBadRequest   = errors.New("bad request")
)

var (
	ErrorUnauthorized = ErrorResponse{
		Status:  "error",
		Message: "Пользователь не авторизован.",
	}

	ErrorBadRequest = ErrorResponse{
		Status:  "error",
		Message: "Ошибка в данных запроса.",
	}

	ErrorInternalServer = ErrorResponse{
		Status:  "error",
		Message: "Ошибка сервера.",
	}
)

var (
	ErrNotFound = errors.New("no record found")
)
