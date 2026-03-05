package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
)

func TestPayload_Valid(t *testing.T) {
	testCases := []struct {
		name          string
		payload       Payload
		expectedError error
	}{
		{
			name: "should return nil when token is not expired",
			payload: Payload{
				ID:        "123",
				Name:      "John Doe",
				Profile:   "CLIENT",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expectedError: nil,
		},
		{
			name: "should return error when token is expired",
			payload: Payload{
				ID:        "123",
				Name:      "John Doe",
				Profile:   "CLIENT",
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expectedError: security.ErrUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.payload.Valid()
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestNewAuthenticator(t *testing.T) {
	auth := NewAuthenticator("secret-key", time.Hour)
	assert.NotNil(t, auth)
}

func TestAuthenticator_Generate(t *testing.T) {
	auth := NewAuthenticator("secret-key", time.Hour)
	user := security.User{
		ID:      "123",
		Name:    "John Doe",
		Profile: "CLIENT",
	}

	token, err := auth.Generate(user)
	assert.NotEmpty(t, token)
	assert.Nil(t, err)
}

func TestAuthenticator_Validate(t *testing.T) {
	auth := NewAuthenticator("secret-key", time.Hour)
	user := security.User{
		ID:      "123",
		Name:    "John Doe",
		Profile: "CLIENT",
	}

	token, _ := auth.Generate(user)

	validatedUser, err := auth.Validate(token)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, validatedUser.ID)
	assert.Equal(t, user.Name, validatedUser.Name)
	assert.Equal(t, user.Profile, validatedUser.Profile)
}

func TestAuthenticator_Validate_InvalidToken(t *testing.T) {
	auth := NewAuthenticator("secret-key", time.Hour)

	_, err := auth.Validate("invalid-token")
	assert.Equal(t, security.ErrUnauthorized, err)
}

func TestAuthenticator_Validate_WrongSecret(t *testing.T) {
	auth1 := NewAuthenticator("secret-key-1", time.Hour)
	auth2 := NewAuthenticator("secret-key-2", time.Hour)
	user := security.User{
		ID:      "123",
		Name:    "John Doe",
		Profile: "CLIENT",
	}

	token, _ := auth1.Generate(user)

	_, err := auth2.Validate(token)
	assert.Equal(t, security.ErrUnauthorized, err)
}
