package handler

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/wellingtonlope/ticket-api/application/usecase"
	"github.com/wellingtonlope/ticket-api/docs"
	"net/http"
)

func InitHandlers(e *echo.Echo, useCases *usecase.AllUseCases) {
	initUserHandler(e, useCases.UserUseCase)
	initTicketHandler(e, useCases.TicketUseCase)

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URL = "/swagger/openapi.yaml"
	}))

	e.GET("/swagger/openapi.yaml", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/yaml", docs.OpenApiYaml)
	})
}

func getAuthorization(c echo.Context) string {
	auths := c.Request().Header["Authorization"]
	if len(auths) > 0 {
		return auths[0]
	}
	return ""
}
