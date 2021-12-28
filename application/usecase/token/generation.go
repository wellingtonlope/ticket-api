package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

func (uc *TokenUseCase) Generate(user *domain.User) (string, *myerrors.Error) {
	claim := UserPayload{
		ID:        user.ID,
		Name:      user.Name,
		Profile:   user.Profile,
		ExpiresAt: time.Now().Add(uc.Duration),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := at.SignedString([]byte(uc.Secret))
	if err != nil {
		return "", myerrors.NewErrorMessage(err.Error(), myerrors.USECASE)
	}

	return token, nil
}
