package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/application/usecase"
)

func InitHandlers(e *echo.Echo, useCases *usecase.AllUseCases) {
	initUserHandler(e, useCases.UserUseCase)
	initTicketHandler(e, useCases.TicketUseCase)
}

func getAuthorization(c echo.Context) string {
	auths := c.Request().Header["Authorization"]
	if len(auths) > 0 {
		return auths[0]
	}
	return ""
}
