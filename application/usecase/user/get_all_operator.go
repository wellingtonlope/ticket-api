package user

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *UserUseCase) GetAllOperator(token string) (*[]domain.User, *myerrors.Error) {
	user, myerr := uc.TokenUseCase.Validate(token)
	if myerr != nil {
		return nil, myerr
	}
	if user.Profile != domain.PROFILE_OPERATOR {
		return nil, myerrors.NewErrorMessage("the user has no permission", myerrors.FORBIDDEN)
	}

	return uc.UserRepository.GetAllOperator()
}
