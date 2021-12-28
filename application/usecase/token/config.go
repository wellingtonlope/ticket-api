package token

import (
	"errors"
	"github.com/wellingtonlope/ticket-api/application/repository"
	"time"
)

type TokenUseCase struct {
	Secret         string
	Duration       time.Duration
	UserRepository repository.UserRepository
}

type UserPayload struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Profile   string    `json:"profile"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (up UserPayload) Valid() error {
	if time.Now().After(up.ExpiresAt) {
		return errors.New("token expired")
	}
	return nil
}
