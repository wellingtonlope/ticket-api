package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllOpen struct {
	ticketRepository repository.TicketRepository
}

func NewGetAllOpen(ticketRepository repository.TicketRepository) *GetAllOpen {
	return &GetAllOpen{ticketRepository: ticketRepository}
}

type GetAllOpenInput struct {
	LoggedUser domain.User
}

type GetAllOpenUserOutput struct {
	ID    string
	Name  string
	Email string
}

type GetAllOpenOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *GetAllOpenUserOutput
	Operator    *GetAllOpenUserOutput
	CreatedAt   *time.Time
	UpdateAt    *time.Time
}

func getAllOpenOutputFromTicket(ticket *domain.Ticket) *GetAllOpenOutput {
	var operator *GetAllOpenUserOutput
	if ticket.Operator != nil {
		operator = &GetAllOpenUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &GetAllOpenOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdateAt:    ticket.UpdatedAt,
		Client: &GetAllOpenUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *GetAllOpen) Handle(input GetAllOpenInput) (*[]GetAllOpenOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, app.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllOpen()
	if err != nil {
		return nil, err
	}

	output := make([]GetAllOpenOutput, 0, len(*tickets))
	for _, ticket := range *tickets {
		output = append(output, *getAllOpenOutputFromTicket(&ticket))
	}

	return &output, nil
}
