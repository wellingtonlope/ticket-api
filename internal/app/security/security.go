package security

import (
	"errors"

	"github.com/wellingtonlope/ticket-api/internal/domain"
)

var (
	ErrForbidden    = errors.New("user don't have permission")
	ErrUnauthorized = errors.New("user is unauthorized")
)

type User struct {
	ID      string
	Name    string
	Profile string
}

func NewUser(user domain.User) User {
	return User{
		ID:      user.ID,
		Name:    user.Name,
		Profile: string(user.Profile),
	}
}

type Authenticator interface {
	Generate(user User) (string, error)
	Validate(token string) (*User, error)
}
