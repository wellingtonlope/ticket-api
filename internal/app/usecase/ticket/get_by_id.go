package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type GetByID struct {
	ticketRepository repository.TicketRepository
}

func NewGetByID(ticketRepository repository.TicketRepository) *GetByID {
	return &GetByID{ticketRepository: ticketRepository}
}

type GetByIDInput struct {
	TicketID   string
	LoggedUser domain.User
}

type GetByIDUserOutput struct {
	ID    string
	Name  string
	Email string
}

type GetByIDOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *GetByIDUserOutput
	Operator    *GetByIDUserOutput
	CreatedAt   *time.Time
	UpdateAt    *time.Time
}

func getByIDoutputFromTicket(ticket *domain.Ticket) *GetByIDOutput {
	var operator *GetByIDUserOutput
	if ticket.Operator != nil {
		operator = &GetByIDUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &GetByIDOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdateAt:    ticket.UpdatedAt,
		Client: &GetByIDUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *GetByID) Handle(input GetByIDInput) (*GetByIDOutput, error) {
	loggedUser := input.LoggedUser
	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket.Client.ID != loggedUser.ID && loggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, repository.ErrTicketNotFound
	}

	return getByIDoutputFromTicket(ticket), nil
}
