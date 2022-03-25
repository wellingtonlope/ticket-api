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
	LoggedUser security.User
}

func (u *AssignToOperator) Handle(input AssignToOperatorInput) (*Output, error) {
	if input.LoggedUser.Profile != string(domain.ProfileOperator) {
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

	return ticketOutputFromTicket(ticket), nil
}
