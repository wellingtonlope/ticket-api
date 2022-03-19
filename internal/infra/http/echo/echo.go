package echo

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/internal/infra/http"
	httpGO "net/http"
)

type Server struct {
	Echo *echo.Echo
}

func (s *Server) Start(port int) error {
	return s.Echo.Start(fmt.Sprintf(":%d", port))
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

		stackMiddleware := r.Handler
		for _, middleware := range r.Middlewares {
			stackMiddleware = middleware(stackMiddleware)
		}

		headers := make(map[string]string)
		headers["Authorization"] = c.Request().Header.Get("Authorization")
		response := stackMiddleware(http.Request{
			Header: headers,
			Body:   string(bytes),
		})

		c.Response().Header().Set("Content-Type", "application/json")
		return c.String(response.HttpCode, response.Body)
	})
}
