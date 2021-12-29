package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TokenUseCase) Validate(token string) (*domain.User, *myerrors.Error) {
	payload := UserPayload{}
	_, err := jwt.ParseWithClaims(token, &payload, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, myerrors.NewErrorMessage("invalid token", myerrors.UNAUTHORIZED)
		}
		return []byte(uc.Secret), nil
	})
	if err != nil {
		return nil, myerrors.NewErrorMessage("invalid token", myerrors.UNAUTHORIZED)
	}

	user, myerr := uc.UserRepository.GetById(payload.ID)
	if myerr != nil {
		if myerr.Type == myerrors.REGISTER_NOT_FOUND {
			return nil, myerrors.NewErrorMessage("invalid token", myerrors.UNAUTHORIZED)
		}
		return nil, myerr
	}

	return user, nil
}
