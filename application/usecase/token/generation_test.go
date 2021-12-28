package token

import (
	"github.com/wellingtonlope/ticket-api/framework/db/local"
	"testing"
	"time"

	"github.com/wellingtonlope/ticket-api/domain"
)

func TestGenerate(t *testing.T) {
	user, _ := domain.UserRegister("name", "email@mail.com", "password")

	t.Run("generate a token", func(t *testing.T) {
		repo := &TokenUseCase{
			Secret:         "123",
			Duration:       time.Minute,
			UserRepository: &local.UserRepositoryLocal{},
		}
		_, _ = repo.UserRepository.Insert(user)

		token, _ := repo.Generate(user)
		if token == "" {
			t.Error("expected a token, but got nothing")
		}

		_, err := repo.Validate(token)
		if err != nil {
			t.Error("expected a token valid, but got an invalid")
		}
	})
}
