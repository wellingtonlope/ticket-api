package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) GetAllOpen(token string) (*[]domain.Ticket, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}
	if user.Profile != domain.PROFILE_OPERATOR {
		return nil, myerrors.NewErrorMessage("the user has no permission", myerrors.FORBIDDEN)
	}

	tickets, myerr := uc.TicketRepository.GetAllOpen()
	if myerr != nil {
		return nil, myerr
	}

	return tickets, nil
}
