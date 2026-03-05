package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

func TestNewUser(t *testing.T) {
	createdAt := time.Now()
	email, _ := domain.NewEmail("john@mail.com")
	password, _ := domain.NewPassword("password123")

	testCases := []struct {
		name           string
		input          domain.User
		expectedResult User
	}{
		{
			name: "should return a new user with operator profile",
			input: domain.User{
				ID:        "123",
				Name:      "John Doe",
				Email:     *email,
				Password:  *password,
				Profile:   domain.ProfileOperator,
				CreatedAt: &createdAt,
			},
			expectedResult: User{
				ID:      "123",
				Name:    "John Doe",
				Profile: "OPERATOR",
			},
		},
		{
			name: "should return a new user with client profile",
			input: domain.User{
				ID:        "456",
				Name:      "Jane Doe",
				Email:     *email,
				Password:  *password,
				Profile:   domain.ProfileClient,
				CreatedAt: &createdAt,
			},
			expectedResult: User{
				ID:      "456",
				Name:    "Jane Doe",
				Profile: "CLIENT",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewUser(tc.input)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
