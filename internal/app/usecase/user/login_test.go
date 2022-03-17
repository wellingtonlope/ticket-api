package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

func TestLogin(t *testing.T) {
	t.Run("should login a user", func(t *testing.T) {
		repo := &memory.UserRepository{}
		usecase := NewLogin(repo)

		email := "user@mail.com"
		password := "password"
		user, _ := domain.UserRegister("user", email, password, time.Now())
		user, _ = repo.Insert(*user)

		input := LoginInput{
			Email:    email,
			Password: password,
		}

		output, err := usecase.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, user.ID, output.ID)
		assert.Equal(t, user.Email.String(), output.Email)
	})

	t.Run("shouldn't login with wrong inputs", func(t *testing.T) {
		repo := &memory.UserRepository{}
		usecase := NewLogin(repo)

		email := "user@mail.com"
		password := "password"
		user, _ := domain.UserRegister("user", email, password, time.Now())
		user, _ = repo.Insert(*user)

		output, err := usecase.Handle(LoginInput{
			Email:    email,
			Password: "wrong password",
		})
		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, ErrUserEmailPasswordWrong, err)

		output, err = usecase.Handle(LoginInput{
			Email:    "wrong@mail.com",
			Password: password,
		})
		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, ErrUserEmailPasswordWrong, err)
	})
}
