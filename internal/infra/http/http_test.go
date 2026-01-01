package http

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
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

func TestWrapError(t *testing.T) {
	t.Run("should return an error in json format", func(t *testing.T) {
		got := wrapError(errors.New("i'm an error"))
		assert.Equal(t, "{\"message\":\"i'm an error\"}", got)
	})

	t.Run("should return an json without message", func(t *testing.T) {
		got := wrapError(nil)
		assert.Equal(t, "{\"message\":\"\"}", got)
	})
}

func TestWrapBody(t *testing.T) {
	t.Run("should return an json with the body", func(t *testing.T) {
		type test struct {
			Foo string `json:"foo"`
		}
		got := wrapBody(test{Foo: "bar"})
		assert.Equal(t, "{\"foo\":\"bar\"}", got)
	})

	t.Run("should return a empty json", func(t *testing.T) {
		got := wrapBody(nil)
		assert.Equal(t, "{}", got)
	})
}
