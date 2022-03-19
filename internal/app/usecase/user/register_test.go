package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	t.Run("should register user", func(t *testing.T) {
		repo := &memory.UserRepository{}
		usecase := NewRegister(repo)

		input := RegisterInput{
			Name:      "user",
			Email:     "user@mail.com",
			Password:  "password",
			CreatedAt: time.Now(),
		}

		output, err := usecase.Handle(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.ID)
		assert.Equal(t, input.Name, output.Name)
		assert.Equal(t, input.Email, output.Email)
		assert.Equal(t, string(domain.PROFILE_CLIENT), output.Profile)
		assert.Equal(t, input.CreatedAt, *output.CreatedAt)
		assert.Nil(t, output.UpdatedAt)

		user, err := repo.GetByID(output.ID)

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, output.ID, user.ID)
		assert.Equal(t, output.Name, user.Name)
		assert.Equal(t, output.Email, user.Email.String())
		assert.Equal(t, output.Profile, string(user.Profile))
		assert.Equal(t, output.CreatedAt, user.CreatedAt)
		assert.Nil(t, user.UpdatedAt)
	})

	t.Run("shouldn't register user when user already register", func(t *testing.T) {
		repo := &memory.UserRepository{}
		usecase := NewRegister(repo)

		input := RegisterInput{
			Name:      "user",
			Email:     "user@mail.com",
			Password:  "password",
			CreatedAt: time.Now(),
		}

		_, err := usecase.Handle(input)
		assert.Nil(t, err)

		output, err := usecase.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, ErrUserAlreadyExists, err)
	})

	t.Run("shouldn't register user when input(email) is invalid", func(t *testing.T) {
		repo := &memory.UserRepository{}
		usecase := NewRegister(repo)

		input := RegisterInput{
			Name:      "user",
			Email:     "user",
			Password:  "password",
			CreatedAt: time.Now(),
		}

		output, err := usecase.Handle(input)

		assert.Nil(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, domain.ErrEmailIsInvalid, err)
	})
}
