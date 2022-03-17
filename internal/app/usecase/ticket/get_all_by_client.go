package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetAllByClient struct {
	ticketRepository repository.TicketRepository
}

func NewGetAllByClient(ticketRepository repository.TicketRepository) *GetAllByClient {
	return &GetAllByClient{ticketRepository: ticketRepository}
}

type GetAllByClientInput struct {
	ClientID   string
	LoggedUser domain.User
}

type GetAllByClientUserOutput struct {
	ID    string
	Name  string
	Email string
}

type GetAllByClientOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *GetAllByClientUserOutput
	Operator    *GetAllByClientUserOutput
	CreatedAt   *time.Time
	UpdateAt    *time.Time
}

func getAllByClientOutputFromTicket(ticket *domain.Ticket) *GetAllByClientOutput {
	var operator *GetAllByClientUserOutput
	if ticket.Operator != nil {
		operator = &GetAllByClientUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &GetAllByClientOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdateAt:    ticket.UpdatedAt,
		Client: &GetAllByClientUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *GetAllByClient) Handle(input GetAllByClientInput) (*[]GetAllByClientOutput, error) {
	if input.ClientID != input.LoggedUser.ID && input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, app.ErrForbidden
	}

	tickets, err := u.ticketRepository.GetAllByClientID(input.ClientID)
	if err != nil {
		return nil, err
	}

	output := make([]GetAllByClientOutput, 0, len(*tickets))
	for _, ticket := range *tickets {
		output = append(output, *getAllByClientOutputFromTicket(&ticket))
	}

	return &output, nil
}
