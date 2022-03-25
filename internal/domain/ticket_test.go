package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenTicket(t *testing.T) {
	title, description, createdAt := "title", "description", time.Now()
	client, _ := UserRegister("client", "client@mail.com", "password", createdAt)

	t.Run("should open a valid ticket", func(t *testing.T) {
		ticket, err := OpenTicket(title, description, createdAt, *client)

		assert.Nil(t, err)
		assert.NotNil(t, ticket)
		assert.Equal(t, title, ticket.Title)
		assert.Equal(t, description, ticket.Description)
		assert.Equal(t, StatusOpen, ticket.Status)
		assert.Equal(t, client.ID, ticket.Client.ID)
		assert.Equal(t, client.Name, ticket.Client.Name)
		assert.Equal(t, client.Email, ticket.Client.Email)
		assert.Equal(t, createdAt, *ticket.CreatedAt)
		assert.Nil(t, ticket.UpdatedAt)
	})

	t.Run("should return an error if title is empty", func(t *testing.T) {
		ticket, err := OpenTicket("", description, createdAt, *client)

		assert.Nil(t, ticket)
		assert.NotNil(t, err)
		assert.Equal(t, ErrTicketTitleIsInvalid, err)
	})
}

func TestTicketGet(t *testing.T) {
	title, description, createdAt, updatedAt := "title", "description", time.Now(), time.Now().Add(time.Hour*48)
	client, _ := UserRegister("client", "client@mail.com", "password", createdAt)
	operator, _ := UserRegister("operator", "operator@mail.com", "password", createdAt)
	operator.Profile = ProfileOperator

	t.Run("should get a valid ticket", func(t *testing.T) {
		ticket, _ := OpenTicket(title, description, createdAt, *client)
		_ = ticket.Get(*operator, updatedAt)

		assert.Equal(t, operator.ID, ticket.Operator.ID)
		assert.Equal(t, operator.Name, ticket.Operator.Name)
		assert.Equal(t, operator.Email, ticket.Operator.Email)
		assert.Equal(t, StatusInProgress, ticket.Status)
		assert.Equal(t, updatedAt, *ticket.UpdatedAt)
	})

	t.Run("should return an error if operator is not an operator", func(t *testing.T) {
		ticket, _ := OpenTicket(title, description, createdAt, *client)
		err := ticket.Get(*client, createdAt)

		assert.NotNil(t, err)
		assert.Equal(t, ErrTicketNoOperator, err)
	})
}

func TestTicketClose(t *testing.T) {
	title, description, createdAt, updatedAt := "title", "description", time.Now(), time.Now().Add(time.Hour*48)
	client, _ := UserRegister("client", "client@mail.com", "password", createdAt)
	operator, _ := UserRegister("operator", "operator@mail.com", "password", createdAt)
	operator.Profile = ProfileOperator

	t.Run("should close a valid ticket", func(t *testing.T) {
		expectedSolution := "solution"
		ticket, _ := OpenTicket(title, description, createdAt, *client)
		_ = ticket.Get(*operator, createdAt)
		err := ticket.Close(expectedSolution, updatedAt)

		assert.Nil(t, err)
		assert.Equal(t, StatusClose, ticket.Status)
		assert.Equal(t, updatedAt, *ticket.UpdatedAt)
		assert.Equal(t, expectedSolution, ticket.Solution)
	})

	t.Run("should return an error if no operator", func(t *testing.T) {
		ticket, _ := OpenTicket(title, description, createdAt, *client)
		err := ticket.Close("solution", createdAt)

		assert.NotNil(t, err)
		assert.Equal(t, ErrTicketNoGetToClose, err)
	})
}
