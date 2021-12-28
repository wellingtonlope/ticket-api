package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
)

func (uc *TicketUseCase) Delete(id, token string) *myerrors.Error {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return myerr
	}

	ticketGet, myerr := uc.TicketRepository.GetById(id)
	if myerr != nil {
		return myerr
	}
	if ticketGet.Client.ID != user.ID {
		return myerrors.NewErrorMessage("you must be the owner of the ticket", myerrors.FORBIDDEN)
	}
	if ticketGet.Operator != nil {
		return myerrors.NewErrorMessage("the ticket must not be picked up from an operator", myerrors.USECASE)
	}

	myerr = uc.TicketRepository.Delete(id)
	if myerr != nil {
		return myerr
	}

	return nil
}
