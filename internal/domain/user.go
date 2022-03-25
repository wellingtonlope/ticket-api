package domain

import (
	"errors"
	"time"
)

var (
	ErrNameIsInvalid = errors.New("name mustn't be empty")
)

type Profile string

const (
	ProfileOperator Profile = "OPERATOR"
	ProfileClient   Profile = "CLIENT"
)

type User struct {
	ID        string
	Name      string
	Email     Email
	Password  Password
	Profile   Profile
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func UserRegister(name, email, password string, createdAt time.Time) (*User, error) {
	if name == "" {
		return nil, ErrNameIsInvalid
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	passwordVO, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:      name,
		Email:     *emailVO,
		Password:  *passwordVO,
		Profile:   ProfileClient,
		CreatedAt: &createdAt,
	}, nil
}
