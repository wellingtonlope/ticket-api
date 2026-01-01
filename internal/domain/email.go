package domain

import (
	"errors"
	"net/mail"
)

var ErrEmailIsInvalid = errors.New("email is invalid")

type Email struct {
	email string
}

func NewEmail(email string) (*Email, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, ErrEmailIsInvalid
	}

	return &Email{email}, nil
}

func (e *Email) String() string {
	return e.email
}
