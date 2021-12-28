package user

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *UserUseCase) Register(name, email, password string) (*UserReponse, *myerrors.Error) {
	userGet, myerr := uc.UserRepository.GetByEmail(email)
	if myerr != nil && myerr.Type != myerrors.REGISTER_NOT_FOUND {
		return nil, myerr
	}

	if userGet != nil {
		return nil, myerrors.NewErrorMessage("email already registered", myerrors.USECASE)
	}

	user, err := domain.UserRegister(name, email, password)
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.DOMAIN)
	}

	userSave, myerr := uc.UserRepository.Insert(user)
	if myerr != nil {
		return nil, myerr
	}

	token, myerr := uc.TokenUseCase.Generate(userSave)
	if myerr != nil {
		return nil, myerr
	}

	return &UserReponse{
		Token: token,
	}, nil
}
