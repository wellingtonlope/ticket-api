package http

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/infra/jwt"
	"net/http"
)

type AuthMiddleware struct {
	authenticator *jwt.Authenticator
}

func NewAuthMiddleware(authenticator *jwt.Authenticator) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: authenticator,
	}
}

func (m *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auths := c.Request().Header["Authorization"]
		if len(auths) > 0 {
			user, err := m.authenticator.Validate(auths[0])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, wrapError(err))
			}
			c.Set(ContextUser, user)
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, wrapError(security.ErrUnauthorized))
	}
}
