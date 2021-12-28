package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) Get(id, token string) (*domain.Ticket, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}

	ticketGet, myerr := uc.TicketRepository.GetById(id)
	if myerr != nil {
		return nil, myerr
	}

	err := ticketGet.Get(user)
	if err != nil {
		if err == domain.ErrTicketNoOperator {
			return nil, myerrors.NewError(err, myerrors.FORBIDDEN)
		}
		return nil, myerrors.NewError(err, myerrors.DOMAIN)
	}

	ticketSaved, myerr := uc.TicketRepository.Update(ticketGet)
	if myerr != nil {
		return nil, myerr
	}

	return ticketSaved, nil
}
