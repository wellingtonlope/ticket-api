package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/ticket"
	"testing"
	"time"
)

func TestTicketUserResponseFromUserOutput(t *testing.T) {
	t.Run("should return ticket user response from user output", func(t *testing.T) {
		output := ticket.UserOutput{ID: "id", Name: "name", Email: "mail@mail.com"}

		got := ticketUserResponseFromUserOutput(&output)
		assert.Equal(t, output.ID, got.ID)
		assert.Equal(t, output.Name, got.Name)
		assert.Equal(t, output.Email, got.Email)
	})
}

func TestTicketResponseFromOutput(t *testing.T) {
	t.Run("should return ticket response from output", func(t *testing.T) {
		dateString := "2019-01-01T00:00:00"
		date, _ := time.Parse("2006-01-02T15:04:05", dateString)
		output := ticket.Output{
			ID:          "id",
			Title:       "title",
			Description: "description",
			Solution:    "solution",
			Status:      "OPEN",
			Client:      &ticket.UserOutput{ID: "id", Name: "name", Email: "mail@mail.com"},
			Operator:    &ticket.UserOutput{ID: "id", Name: "name", Email: "mail@mail.com"},
			CreatedAt:   &date,
			UpdatedAt:   &date,
		}

		got := ticketResponseFromOutput(output)

		assert.Equal(t, output.ID, got.ID)
		assert.Equal(t, output.Title, got.Title)
		assert.Equal(t, output.Description, got.Description)
		assert.Equal(t, output.Solution, got.Solution)
		assert.Equal(t, output.Status, got.Status)
		assert.Equal(t, output.Client.ID, got.Client.ID)
		assert.Equal(t, output.Client.Name, got.Client.Name)
		assert.Equal(t, output.Client.Email, got.Client.Email)
		assert.Equal(t, output.Operator.ID, got.Operator.ID)
		assert.Equal(t, output.Operator.Name, got.Operator.Name)
		assert.Equal(t, output.Operator.Email, got.Operator.Email)
		assert.Equal(t, dateString, got.CreatedAt)
		assert.Equal(t, dateString, got.UpdatedAt)
	})

	t.Run("should return an empty ticket response", func(t *testing.T) {
		got := ticketResponseFromOutput(ticket.Output{})

		assert.Equal(t, "", got.ID)
		assert.Equal(t, "", got.Title)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, "", got.Solution)
		assert.Equal(t, "", got.Status)
		assert.Nil(t, got.Client)
		assert.Nil(t, got.Operator)
		assert.Equal(t, "", got.CreatedAt)
		assert.Equal(t, "", got.UpdatedAt)
	})
}

func TestTicketsResponseFromOutputs(t *testing.T) {
	t.Run("should return tickets response from outputs", func(t *testing.T) {
		got := ticketsResponseFromOutputs([]ticket.Output{{}, {}})

		assert.Len(t, *got, 2)
	})
}
