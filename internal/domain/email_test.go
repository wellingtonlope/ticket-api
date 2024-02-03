package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedResult *Email
		expectedError  error
	}{
		{
			name:           "should return a new email",
			input:          "email@mail.com",
			expectedResult: &Email{email: "email@mail.com"},
			expectedError:  nil,
		},
		{
			name:           "should return an error when empty email",
			input:          "",
			expectedResult: nil,
			expectedError:  ErrEmailIsInvalid,
		},
		{
			name:           "should return an error when invalid email",
			input:          "email",
			expectedResult: nil,
			expectedError:  ErrEmailIsInvalid,
		},
		{
			name:           "should return a email when valid email .br",
			input:          "email@mail.com.br",
			expectedResult: &Email{email: "email@mail.com.br"},
			expectedError:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewEmail(tc.input)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestEmail_String(t *testing.T) {
	expectedString := "email@mail.com"
	got, _ := NewEmail(expectedString)
	assert.Equal(t, expectedString, got.String())
}
