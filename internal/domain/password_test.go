package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword(t *testing.T) {
	t.Run("should return a new password", func(t *testing.T) {
		expectedPassword := "password"
		got, err := NewPassword(expectedPassword)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.True(t, got.IsCorrectPassword(expectedPassword))
	})

	t.Run("should return an error when empty password", func(t *testing.T) {
		got, err := NewPassword("")

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, ErrPasswordIsInvalid, err)
	})

	t.Run("should return an error when invalid password", func(t *testing.T) {
		got, err := NewPassword("123")

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, ErrPasswordIsInvalid, err)
	})
}

func TestIsCorrectPassword(t *testing.T) {
	t.Run("should return true when password is correct", func(t *testing.T) {
		expectedPassword := "password"
		got, err := NewPassword(expectedPassword)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.True(t, got.IsCorrectPassword(expectedPassword))
	})

	t.Run("should return false when password is incorrect", func(t *testing.T) {
		expectedPassword := "password"
		got, err := NewPassword(expectedPassword)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.False(t, got.IsCorrectPassword("incorrect"))
	})
}

func TestNewPasswordHashed(t *testing.T) {
	t.Run("should return a new password with hashed string password", func(t *testing.T) {
		expectedPassword := "password"
		got, err := NewPasswordHashed(expectedPassword)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, expectedPassword, got.String())
	})
}
