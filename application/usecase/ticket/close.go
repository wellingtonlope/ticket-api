package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) Close(id, solution, token string) (*domain.Ticket, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}
	if user.Profile != domain.PROFILE_OPERATOR {
		return nil, myerrors.NewErrorMessage("the user has no permission", myerrors.FORBIDDEN)
	}

	ticketGet, myerr := uc.TicketRepository.GetById(id)
	if myerr != nil {
		return nil, myerr
	}

	err := ticketGet.Close(solution)
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.DOMAIN)
	}

	ticketSaved, myerr := uc.TicketRepository.Update(ticketGet)
	if myerr != nil {
		return nil, myerr
	}

	return ticketSaved, nil
}
