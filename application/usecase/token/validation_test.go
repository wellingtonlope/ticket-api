package token

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/framework/db/local"
	"testing"
	"time"

	"github.com/wellingtonlope/ticket-api/domain"
)

func TestIsValidToken(t *testing.T) {
	user, _ := domain.UserRegister("name", "email@mail.com", "password")
	uc := &TokenUseCase{
		Secret:         "123",
		Duration:       time.Minute,
		UserRepository: &local.UserRepositoryLocal{}}
	_, _ = uc.UserRepository.Insert(user)
	validToken, _ := uc.Generate(user)

	t.Run("a valid token", func(t *testing.T) {
		userPayload, _ := uc.Validate(validToken)
		if userPayload == nil {
			t.Error("expected an user payload, but got a nil")
		}
	})

	t.Run("an invalid token", func(t *testing.T) {
		_, err := uc.Validate("123")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		ucExpired := &TokenUseCase{
			Secret:         "123",
			Duration:       time.Minute * -1,
			UserRepository: &local.UserRepositoryLocal{},
		}
		expiredToken, _ := ucExpired.Generate(user)
		_, err := ucExpired.Validate(expiredToken)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("Non-existent user", func(t *testing.T) {
		_, err := uc.Validate("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjQ2MmI0OGE1LWYyMzEtNGVhNS05OTgwLWFkZDJkZGY1YjViNCIsIm5hbWUiOiJXZWxsaW5ndG9uIiwicHJvZmlsZSI6IkNMSUVOVCIsImV4cGlyZXNBdCI6IjIwMjEtMTItMzBUMjA6MTg6MzUuNTEzNjEzOTI5LTAzOjAwIn0.5NI3DForgk3Zy_MiSBbPUkhkpH9GF3c4jz7AMyHYMvo")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
		if err.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, err.Type)
		}
	})
}
