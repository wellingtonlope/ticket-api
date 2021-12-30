package user

import (
	"testing"

	"github.com/wellingtonlope/ticket-api/application/myerrors"
)

func TestRegister(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		uc := newUserUC()
		got, _ := uc.Register("client", "client@mail.com", "password")
		if got == nil {
			t.Error("expected a userReponse, but got a nil")
		}
		if got != nil && got.Token == "" {
			t.Error("expected a token, but got nothing")
		}
	})

	t.Run("invalid user", func(t *testing.T) {
		uc := newUserUC()
		_, err := uc.Register("client", "client", "password")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
		if err != nil && err.Type != myerrors.DOMAIN {
			t.Errorf("expected %q, but got %q", myerrors.DOMAIN, err.Type)
		}
	})

	t.Run("duplicated email", func(t *testing.T) {
		uc := newUserUC()
		email := "client@mail.com"
		uc.Register("client", email, "password")
		_, err := uc.Register("client", email, "password")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}
		if err != nil && err.Type != myerrors.USECASE {
			t.Errorf("expected %q, but got %q", myerrors.USECASE, err.Type)
		}
	})
}
