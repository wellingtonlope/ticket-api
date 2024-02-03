package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		assertResult  func(assert.TestingT, interface{}, ...interface{}) bool
		expectedError error
	}{
		{
			name:          "should return a new password",
			input:         "password",
			assertResult:  assert.NotNil,
			expectedError: nil,
		},
		{
			name:          "should return an error when empty password",
			input:         "",
			assertResult:  assert.Nil,
			expectedError: ErrPasswordIsInvalid,
		},
		{
			name:          "should return an error when invalid password",
			input:         "123",
			assertResult:  assert.Nil,
			expectedError: ErrPasswordIsInvalid,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewPassword(tc.input)
			tc.assertResult(t, got)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestIsCorrectPassword(t *testing.T) {
	testCases := []struct {
		name           string
		password       Password
		input          string
		expectedResult bool
	}{
		{
			name: "should return true when password is correct",
			password: func() Password {
				m, _ := NewPassword("password")
				return *m
			}(),
			input:          "password",
			expectedResult: true,
		},
		{
			name: "should return false when password is incorrect",
			password: func() Password {
				m, _ := NewPassword("password")
				return *m
			}(),
			input:          "incorrect",
			expectedResult: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.password.IsCorrectPassword(tc.input)
			assert.Equal(t, tc.expectedResult, got)
		})
	}
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
