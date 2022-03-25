package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type UserOutput struct {
	ID    string
	Name  string
	Email string
}

type Output struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *UserOutput
	Operator    *UserOutput
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func ticketOutputFromTicket(ticket *domain.Ticket) *Output {
	var operator *UserOutput
	if ticket.Operator != nil {
		operator = &UserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &Output{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Client: &UserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func ticketsOutputsFromTickets(tickets *[]domain.Ticket) *[]Output {
	outputs := make([]Output, 0, len(*tickets))
	for _, ticket := range *tickets {
		outputs = append(outputs, *ticketOutputFromTicket(&ticket))
	}
	return &outputs
}
