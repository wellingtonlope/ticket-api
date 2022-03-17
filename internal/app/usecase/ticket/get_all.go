package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAll struct {
	ticketRepository repository.TicketRepository
}

func NewGetAll(ticketRepository repository.TicketRepository) *GetAll {
	return &GetAll{ticketRepository: ticketRepository}
}

type GetAllInput struct {
	LoggedUser domain.User
}

type GetAllUserOutput struct {
	ID    string
	Name  string
	Email string
}

type GetAllOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *GetAllUserOutput
	Operator    *GetAllUserOutput
	CreatedAt   *time.Time
	UpdateAt    *time.Time
}

func getAllOutputFromTicket(ticket *domain.Ticket) *GetAllOutput {
	var operator *GetAllUserOutput
	if ticket.Operator != nil {
		operator = &GetAllUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &GetAllOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdateAt:    ticket.UpdatedAt,
		Client: &GetAllUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *GetAll) Handle(input GetAllInput) (*[]GetAllOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, app.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAll()
	if err != nil {
		return nil, err
	}

	output := make([]GetAllOutput, 0, len(*tickets))
	for _, ticket := range *tickets {
		output = append(output, *getAllOutputFromTicket(&ticket))
	}

	return &output, nil
}
