package user

import "github.com/wellingtonlope/ticket-api/application/myerrors"

func (uc *UserUseCase) Login(email, password string) (*UserReponse, *myerrors.Error) {
	userGet, myerr := uc.UserRepository.GetByEmail(email)
	if myerr != nil {
		if myerr.Type == myerrors.REGISTER_NOT_FOUND {
			return nil, myerrors.NewErrorMessage("user and password are wrong", myerrors.UNAUTHORIZED)
		}

		return nil, myerr
	}

	isCorrectPassword := userGet.IsCorrectPassword(password)
	if !isCorrectPassword {
		return nil, myerrors.NewErrorMessage("user and password are wrong", myerrors.UNAUTHORIZED)
	}

	token, myerr := uc.TokenUseCase.Generate(userGet)
	if myerr != nil {
		return nil, myerr
	}

	return &UserReponse{
		Token: token,
	}, nil
}
