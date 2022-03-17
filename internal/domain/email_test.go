package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	t.Run("should return a new email", func(t *testing.T) {
		expectedEmail := "email@mail.com"
		got, err := NewEmail(expectedEmail)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, expectedEmail, got.String())
	})

	t.Run("should return an error when empty email", func(t *testing.T) {
		got, err := NewEmail("")

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmailIsInvalid, err)
	})

	t.Run("should return an error when invalid email", func(t *testing.T) {
		got, err := NewEmail("email")

		assert.Nil(t, got)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmailIsInvalid, err)
	})

	t.Run("should return a email when valid email .br", func(t *testing.T) {
		expectedEmail := "email@mail.com.br"
		got, err := NewEmail(expectedEmail)

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, expectedEmail, got.String())
	})
}
