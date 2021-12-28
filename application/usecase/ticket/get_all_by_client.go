package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) GetAllByClient(idClient, token string) (*[]domain.Ticket, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}
	if idClient != user.ID && user.Profile != domain.PROFILE_OPERATOR {
		return nil, myerrors.NewErrorMessage("the user has no permission", myerrors.FORBIDDEN)
	}

	client, myerr := uc.UserUseCase.GetById(idClient)
	if myerr != nil {
		return nil, myerr
	}

	tickets, myerr := uc.TicketRepository.GetAllByClient(client)
	if myerr != nil {
		return nil, myerr
	}

	return tickets, nil
}
