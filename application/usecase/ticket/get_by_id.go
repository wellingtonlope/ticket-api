package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) GetById(id, token string) (*domain.Ticket, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}

	ticket, myerr := uc.TicketRepository.GetById(id)
	if myerr != nil {
		return nil, myerr
	}
	if ticket.Client.ID != user.ID && user.Profile != domain.PROFILE_OPERATOR {
		return nil, myerrors.NewErrorMessage("the user has no permission", myerrors.FORBIDDEN)
	}

	return ticket, nil
}
