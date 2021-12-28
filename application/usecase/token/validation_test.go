package token

import (
	"github.com/wellingtonlope/ticket-api/framework/db/local"
	"testing"
	"time"

	"github.com/wellingtonlope/ticket-api/domain"
)

func TestIsValidToken(t *testing.T) {
	user, _ := domain.UserRegister("name", "email@mail.com", "password")
	repo := &TokenUseCase{
		Secret:         "123",
		Duration:       time.Minute,
		UserRepository: &local.UserRepositoryLocal{}}
	_, _ = repo.UserRepository.Insert(user)
	validToken, _ := repo.Generate(user)

	t.Run("a valid token", func(t *testing.T) {
		userPayload, _ := repo.Validate(validToken)
		if userPayload == nil {
			t.Error("expected an user payload, but got a nil")
		}
	})

	t.Run("an invalid token", func(t *testing.T) {
		_, err := repo.Validate("123")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		repoExpired := &TokenUseCase{
			Secret:         "123",
			Duration:       time.Minute * -1,
			UserRepository: &local.UserRepositoryLocal{},
		}
		expiredToken, _ := repoExpired.Generate(user)
		_, err := repoExpired.Validate(expiredToken)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
	})
}
