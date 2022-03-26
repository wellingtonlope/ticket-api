package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/user"
	"testing"
	"time"
)

func TestUserResponseFromUserOutput(t *testing.T) {
	t.Run("should return user response from user output", func(t *testing.T) {
		dateString := "2019-01-01T00:00:00"
		date, _ := time.Parse("2006-01-02T15:04:05", dateString)
		output := user.Output{
			ID:        "id",
			Name:      "name",
			Email:     "mail@mail.com",
			Profile:   "CLIENT",
			CreatedAt: &date,
			UpdatedAt: &date,
		}
		response := userResponseFromUserOutput(output)
		assert.Equal(t, output.ID, response.ID)
		assert.Equal(t, output.Name, response.Name)
		assert.Equal(t, output.Email, response.Email)
		assert.Equal(t, dateString, response.CreatedAt)
		assert.Equal(t, dateString, response.UpdatedAt)
	})
	t.Run("should return an empty user response", func(t *testing.T) {
		response := userResponseFromUserOutput(user.Output{})
		assert.Equal(t, "", response.ID)
		assert.Equal(t, "", response.Name)
		assert.Equal(t, "", response.Email)
		assert.Equal(t, "", response.CreatedAt)
		assert.Equal(t, "", response.UpdatedAt)
	})
}

func TestUsersResponseFromUsersOutput(t *testing.T) {
	t.Run("should return users response from users output", func(t *testing.T) {
		response := usersResponseFromUsersOutput([]user.Output{{}, {}})
		assert.Len(t, *response, 2)
	})
}
