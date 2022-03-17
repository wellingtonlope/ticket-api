package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Get struct {
	ticketRepository repository.TicketRepository
}

func NewGet(ticketRepository repository.TicketRepository) *Get {
	return &Get{ticketRepository: ticketRepository}
}

type GetInput struct {
	TicketID   string
	UpdatedAt  time.Time
	LoggedUser domain.User
}

type GetUserOutput struct {
	ID    string
	Name  string
	Email string
}

type GetOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *GetUserOutput
	Operator    *GetUserOutput
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func getOutputFromTicket(ticket *domain.Ticket) *GetOutput {
	var operator *GetUserOutput
	if ticket.Operator != nil {
		operator = &GetUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &GetOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Client: &GetUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *Get) Handle(input GetInput) (*GetOutput, error) {
	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	err = ticket.Get(input.LoggedUser, input.UpdatedAt)
	if err != nil {
		if err == domain.ErrTicketNoOperator {
			return nil, security.ErrForbidden
		}
		return nil, err
	}

	ticket, err = u.ticketRepository.Update(*ticket)
	if err != nil {
		return nil, err
	}

	return getOutputFromTicket(ticket), nil
}
