package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) AssingToOperator(idTicket, idOperator, token string) (*domain.Ticket, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}
	if user.Profile != domain.PROFILE_OPERATOR {
		return nil, myerrors.NewErrorMessage("the user has no permission", myerrors.FORBIDDEN)
	}

	operator, myerr := uc.UserUseCase.GetById(idOperator)
	if myerr != nil {
		return nil, myerr
	}

	ticketGet, myerr := uc.TicketRepository.GetById(idTicket)
	if myerr != nil {
		return nil, myerr
	}

	err := ticketGet.Get(operator)
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.DOMAIN)
	}

	ticketSaved, myerr := uc.TicketRepository.Update(ticketGet)
	if myerr != nil {
		return nil, myerr
	}

	return ticketSaved, nil
}
