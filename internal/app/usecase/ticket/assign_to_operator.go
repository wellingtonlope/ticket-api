package ticket

import (
	"time"

	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type AssignToOperator struct {
	ticketRepository repository.TicketRepository
	userRepository   repository.UserRepository
}

func NewAssignToOperator(ticketRepository repository.TicketRepository, userRepository repository.UserRepository) *AssignToOperator {
	return &AssignToOperator{ticketRepository: ticketRepository, userRepository: userRepository}
}

type AssignToOperatorInput struct {
	TicketID   string
	OperatorID string
	UpdatedAt  time.Time
	LoggedUser domain.User
}

type AssignToOperatorUserOutput struct {
	ID    string
	Name  string
	Email string
}

type AssignToOperatorOutput struct {
	ID          string
	Title       string
	Description string
	Solution    string
	Status      string
	Client      *AssignToOperatorUserOutput
	Operator    *AssignToOperatorUserOutput
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func assignToOperatorOutputFromTicket(ticket *domain.Ticket) *AssignToOperatorOutput {
	var operator *AssignToOperatorUserOutput
	if ticket.Operator != nil {
		operator = &AssignToOperatorUserOutput{
			ID:    ticket.Operator.ID,
			Name:  ticket.Operator.Name,
			Email: ticket.Operator.Email.String(),
		}
	}

	return &AssignToOperatorOutput{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Solution:    ticket.Solution,
		Status:      string(ticket.Status),
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Client: &AssignToOperatorUserOutput{
			ID:    ticket.Client.ID,
			Name:  ticket.Client.Name,
			Email: ticket.Client.Email.String(),
		},
		Operator: operator,
	}
}

func (u *AssignToOperator) Handle(input AssignToOperatorInput) (*AssignToOperatorOutput, error) {
	if input.LoggedUser.Profile != domain.PROFILE_OPERATOR {
		return nil, security.ErrForbidden
	}

	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetByID(input.OperatorID)
	if err != nil {
		return nil, err
	}

	err = ticket.Get(*user, input.UpdatedAt)
	if err != nil {
		return nil, err
	}

	ticket, err = u.ticketRepository.Update(*ticket)
	if err != nil {
		return nil, err
	}

	return assignToOperatorOutputFromTicket(ticket), nil
}
