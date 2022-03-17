package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllByOperator struct {
	ticketRepository repository.TicketRepository
}

func NewGetAllByOperator(ticketRepository repository.TicketRepository) *GetAllByOperator {
	return &GetAllByOperator{ticketRepository: ticketRepository}
}

type GetAllByOperatorInput struct {
	OperatorID string
	LoggedUser domain.User
}

type GetAllByOperatorUserOutput struct {
	ID    string
	Name  string
	Email string
}

type GetAllByOperatorOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *GetAllByOperatorUserOutput
	Operator    *GetAllByOperatorUserOutput
	CreatedAt   *time.Time
	UpdateAt    *time.Time
}

func getAllByOperatorOutputFromTicket(ticket *domain.Ticket) *GetAllByOperatorOutput {
	var operator *GetAllByOperatorUserOutput
	if ticket.Operator != nil {
		operator = &GetAllByOperatorUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &GetAllByOperatorOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdateAt:    ticket.UpdatedAt,
		Client: &GetAllByOperatorUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *GetAllByOperator) Handle(input GetAllByOperatorInput) (*[]GetAllByOperatorOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, app.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllByOperatorID(input.OperatorID)
	if err != nil {
		return nil, err
	}

	output := make([]GetAllByOperatorOutput, 0, len(*tickets))
	for _, ticket := range *tickets {
		output = append(output, *getAllByOperatorOutputFromTicket(&ticket))
	}

	return &output, nil
}
