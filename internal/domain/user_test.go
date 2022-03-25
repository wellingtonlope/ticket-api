package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	name, email, password, createdAt := "name", "email@mail.com", "password", time.Now()

	t.Run("should register a user", func(t *testing.T) {
		user, err := UserRegister(name, email, password, createdAt)

		assert.Nil(t, err)
		assert.Equal(t, name, user.Name)
		assert.Equal(t, email, user.Email.String())
		assert.True(t, user.Password.IsCorrectPassword(password))
		assert.Equal(t, ProfileClient, user.Profile)
		assert.Equal(t, createdAt, *user.CreatedAt)
	})

	t.Run("shouldn't register a user when empty name", func(t *testing.T) {
		user, err := UserRegister("", email, password, createdAt)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNameIsInvalid, err)
	})

	t.Run("shouldn't register a user when invalid email", func(t *testing.T) {
		user, err := UserRegister(name, "email", password, createdAt)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmailIsInvalid, err)
	})

	t.Run("shouldn't register a user when invalid password", func(t *testing.T) {
		user, err := UserRegister(name, email, "", createdAt)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, ErrPasswordIsInvalid, err)
	})
}
