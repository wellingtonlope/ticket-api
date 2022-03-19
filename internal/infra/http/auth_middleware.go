package http

import (
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/infra/jwt"
	"net/http"
)

type AuthMiddleware struct {
	Authenticator *jwt.Authenticator
}

func (m *AuthMiddleware) Handle(next Handler) Handler {
	return func(r Request) Response {
		jwtToken := r.Header["Authorization"]

		if len(jwtToken) > 0 {
			user, err := m.Authenticator.Validate(jwtToken)
			if err != nil {
				return Response{
					HttpCode: http.StatusUnauthorized,
					Body:     wrapError(err),
				}
			}
			r.LoggedUser = user
			return next(r)
		}

		return Response{
			HttpCode: http.StatusUnauthorized,
			Body:     wrapError(security.ErrUnauthorized),
		}
	}
}
