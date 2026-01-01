package echo

import (
	"encoding/json"
	"fmt"
	httpGO "net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/wellingtonlope/ticket-api/internal/infra/http"
)

type Server struct {
	Echo *echo.Echo
}

func (s *Server) Start(port int) error {
	return s.Echo.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) RegisterSwagger(file []byte) {
	s.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{"/swagger/openapi.yaml"}
	}))
	s.Echo.GET("/swagger/openapi.yaml", func(c echo.Context) error {
		return c.Blob(httpGO.StatusOK, "application/yaml", file)
	})
}

func (s *Server) Register(r http.Route) {
	s.Echo.Add(r.Method, r.Path, func(c echo.Context) error {
		var request interface{}
		err := c.Bind(&request)
		if err != nil {
			return c.JSON(httpGO.StatusBadRequest, http.ErrorResponse{
				Message: err.Error(),
			})
		}
		bytes, err := json.Marshal(request)
		if err != nil {
			return c.JSON(httpGO.StatusBadRequest, http.ErrorResponse{
				Message: err.Error(),
			})
		}

		params := make(map[string]string)
		for _, key := range c.ParamNames() {
			params[key] = c.Param(key)
		}

		stackMiddleware := r.Handler
		for _, middleware := range r.Middlewares {
			stackMiddleware = middleware(stackMiddleware)
		}

		headers := make(map[string]string)
		headers[http.AuthorizationHeader] = c.Request().Header.Get(http.AuthorizationHeader)
		response := stackMiddleware(http.Request{
			Header: headers,
			Params: params,
			Body:   string(bytes),
		})

		c.Response().Header().Set(http.ContentTypeHeader, http.ContentTypeJSON)
		return c.String(response.HttpCode, response.Body)
	})
}
