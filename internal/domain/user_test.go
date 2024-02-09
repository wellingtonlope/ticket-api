package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	exampleDate := time.Now()
	type args struct {
		name      string
		email     string
		password  string
		createdAt time.Time
	}
	testCases := []struct {
		name          string
		args          args
		assertResult  func(t *testing.T, got *User)
		expectedError error
	}{
		{
			name: "should register a user",
			args: args{
				name:      "name",
				email:     "email@mail.com",
				password:  "password",
				createdAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *User) {
				assert.Equal(t, "name", got.Name)
				assert.Equal(t, "email@mail.com", got.Email.String())
				assert.True(t, got.Password.IsCorrectPassword("password"))
				assert.Equal(t, ProfileClient, got.Profile)
				assert.Equal(t, exampleDate, *got.CreatedAt)
			},
			expectedError: nil,
		},
		{
			name: "shouldn't register a user when empty name",
			args: args{
				name:      "",
				email:     "email@mail.com",
				password:  "password",
				createdAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *User) {
				assert.Nil(t, got)
			},
			expectedError: ErrNameIsInvalid,
		},
		{
			name: "shouldn't register a user when invalid email",
			args: args{
				name:      "name",
				email:     "email",
				password:  "password",
				createdAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *User) {
				assert.Nil(t, got)
			},
			expectedError: ErrEmailIsInvalid,
		},
		{
			name: "shouldn't register a user when invalid password",
			args: args{
				name:      "name",
				email:     "email@mail.com",
				password:  "",
				createdAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *User) {
				assert.Nil(t, got)
			},
			expectedError: ErrPasswordIsInvalid,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := UserRegister(tc.args.name, tc.args.email, tc.args.password, tc.args.createdAt)
			tc.assertResult(t, got)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
