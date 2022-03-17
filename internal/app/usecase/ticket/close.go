package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Close struct {
	ticketRepository repository.TicketRepository
}

func NewClose(ticketRepository repository.TicketRepository) *Close {
	return &Close{ticketRepository: ticketRepository}
}

type CloseInput struct {
	TicketID   string
	Solution   string
	UpdatedAt  time.Time
	LoggedUser domain.User
}

type CloseUserOutput struct {
	ID    string
	Name  string
	Email string
}

type CloseOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *CloseUserOutput
	Operator    *CloseUserOutput
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func closeOutputFromTicket(ticket *domain.Ticket) *CloseOutput {
	var operator *CloseUserOutput
	if ticket.Operator != nil {
		operator = &CloseUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &CloseOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Client: &CloseUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *Close) Handle(input CloseInput) (*CloseOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, app.ErrForbidden
	}

	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	err = ticket.Close(input.Solution, input.UpdatedAt)
	if err != nil {
		return nil, err
	}

	ticket, err = u.ticketRepository.Update(*ticket)
	if err != nil {
		return nil, err
	}

	return closeOutputFromTicket(ticket), nil
}
