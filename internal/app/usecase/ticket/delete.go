package ticket

import (
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type Delete interface {
	Handle(input DeleteInput) (*DeleteOutput, error)
}

type delete struct {
	ticketRepository repository.TicketRepository
}

func NewDelete(ticketRepository repository.TicketRepository) Delete {
	return &delete{ticketRepository: ticketRepository}
}

type DeleteInput struct {
	TicketID   string
	LoggedUser security.User
}

type DeleteOutput struct{}

func (u *delete) Handle(input DeleteInput) (*DeleteOutput, error) {
	ticket, err := u.ticketRepository.GetByID(input.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket.Client.ID != input.LoggedUser.ID {
		return nil, security.ErrForbidden
	}

	if ticket.Status != domain.StatusOpen {
		return nil, security.ErrForbidden
	}

	err = u.ticketRepository.DeleteByID(ticket.ID)
	if err != nil {
		return nil, err
	}

	return &DeleteOutput{}, nil
}
