package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordIsInvalid = errors.New("password must be 6 characters or longer")

type Password struct {
	hashedPassword string
}

func NewPassword(password string) (*Password, error) {
	if len(password) < 6 {
		return nil, ErrPasswordIsInvalid
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Password{string(hashPassword)}, nil
}

func NewPasswordHashed(hashedPassword string) (*Password, error) {
	return &Password{hashedPassword}, nil
}

func (p *Password) String() string {
	return p.hashedPassword
}

func (p *Password) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedPassword), []byte(password))
	return err == nil
}
