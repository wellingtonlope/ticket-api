package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) Open(title, description, token string) (*domain.Ticket, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}

	ticket, err := domain.TicketOpen(title, description, user)
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.DOMAIN)
	}

	ticketSaved, myerr := uc.TicketRepository.Insert(ticket)
	if myerr != nil {
		return nil, myerr
	}

	return ticketSaved, nil
}
