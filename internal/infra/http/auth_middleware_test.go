package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	httpGo "net/http"
	"testing"
)

type mockAuth struct {
	mock.Mock
}

func (m *mockAuth) Generate(user security.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *mockAuth) Validate(token string) (*security.User, error) {
	args := m.Called(token)
	var result *security.User
	if args.Get(0) != nil {
		result = args.Get(0).(*security.User)
	}
	return result, args.Error(1)
}

func TestAuthMiddleware_Handle(t *testing.T) {
	tokens := []string{"", "token"}
	for _, token := range tokens {
		t.Run("should return error when token is empty or invalid", func(t *testing.T) {
			mockAuth := new(mockAuth)
			mockAuth.On("Validate", mock.Anything).Return(nil, security.ErrUnauthorized)

			middleware := AuthMiddleware{
				Authenticator: mockAuth,
			}

			response := middleware.Handle(func(_ Request) Response {
				return Response{
					HttpCode: httpGo.StatusOK,
					Body:     "{\"message\": \"ok\"}",
				}
			})(Request{
				Header: map[string]string{
					"Authorization": token,
				},
			})

			assert.Equal(t, httpGo.StatusUnauthorized, response.HttpCode)
			assert.Equal(t, "{\"message\":\"user is unauthorized\"}", response.Body)
		})
	}

	t.Run("should pass with a logged user when valid token", func(t *testing.T) {
		securityUser := &security.User{
			ID:      "id",
			Name:    "test",
			Profile: "CLIENT",
		}
		mockAuth := new(mockAuth)
		mockAuth.On("Validate", mock.Anything).Return(securityUser, nil)

		middleware := AuthMiddleware{
			Authenticator: mockAuth,
		}

		response := middleware.Handle(func(r Request) Response {
			assert.Equal(t, securityUser, r.LoggedUser)

			return Response{
				HttpCode: httpGo.StatusOK,
				Body:     "{\"message\": \"ok\"}",
			}
		})(Request{
			Header: map[string]string{
				"Authorization": "Bearer 123",
			},
		})

		assert.Equal(t, httpGo.StatusOK, response.HttpCode)
		assert.Equal(t, "{\"message\": \"ok\"}", response.Body)
	})
}
