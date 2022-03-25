package http

import (
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"net/http"
)

type AuthMiddleware struct {
	Authenticator security.Authenticator
}

func (m *AuthMiddleware) Handle(next Handler) Handler {
	return func(r Request) Response {
		jwtToken := r.Header[AuthorizationHeader]

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
