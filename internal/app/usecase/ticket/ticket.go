package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type TicketUserOutput struct {
	ID    string
	Name  string
	Email string
}

type TicketOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *TicketUserOutput
	Operator    *TicketUserOutput
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func ticketOutputFromTicket(ticket *domain.Ticket) *TicketOutput {
	var operator *TicketUserOutput
	if ticket.Operator != nil {
		operator = &TicketUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &TicketOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Client: &TicketUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func ticketsOutputsFromTickets(tickets *[]domain.Ticket) *[]TicketOutput {
	outputs := make([]TicketOutput, 0, len(*tickets))
	for _, ticket := range *tickets {
		outputs = append(outputs, *ticketOutputFromTicket(&ticket))
	}
	return &outputs
}
