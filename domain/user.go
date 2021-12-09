package domain

import (
	"errors"
	"net/mail"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	PROFILE_OPERATOR = "OPERATOR"
	PROFILE_CLIENT   = "CLIENT"
)

type User struct {
	Base     `bson:",inline"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Profile  string `json:"profile"`
}

func UserRegister(name, email, password string) (*User, error) {
	if name == "" {
		return nil, errors.New("name musn't be empty")
	}

	if email == "" {
		return nil, errors.New("email musn't be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, errors.New("email is invalid")
	}

	if password == "" {
		return nil, errors.New("password musn't be empty")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be 6 characters or longer")
	}

	now := time.Now()
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error during generate a password with bcrypt")
	}

	return &User{
		Base: Base{
			ID:        uuid.NewV4().String(),
			CreatedAt: &now,
		},
		Name:     name,
		Email:    email,
		Password: string(hashPassword),
		Profile:  PROFILE_CLIENT,
	}, nil
}

func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
