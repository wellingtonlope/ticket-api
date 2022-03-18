package http

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase"
	"github.com/wellingtonlope/ticket-api/internal/infra/jwt"
)

const (
	ContextUser = "LOGGED_USER"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Http struct {
	Echo          *echo.Echo
	UseCases      *usecase.AllUseCases
	Repositories  *repository.AllRepositories
	Authenticator *jwt.Authenticator
}

func wrapError(err error) *ErrorResponse {
	if err == nil {
		return nil
	}

	return &ErrorResponse{
		Message: err.Error(),
	}
}

func (http Http) Init() {
	initUserHandler(http.Echo, http.UseCases, http.Authenticator)
}
