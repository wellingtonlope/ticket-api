package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TicketUseCase) GetAllByOperator(idOperator, token string) (*[]domain.Ticket, *myerrors.Error) {
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

	tickets, myerr := uc.TicketRepository.GetAllByOperator(operator)
	if myerr != nil {
		return nil, myerr
	}

	return tickets, nil
}
