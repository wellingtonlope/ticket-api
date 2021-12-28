package user

import (
	"testing"

	"github.com/wellingtonlope/ticket-api/application/myerrors"
)

func TestLogin(t *testing.T) {
	uc := newUserUC()
	email, password := "client@mail.com", "password"
	uc.Register("client", email, password)

	t.Run("a valid user", func(t *testing.T) {
		got, _ := uc.Login(email, password)
		if got == nil {
			t.Error("expected an UserResponse, but got a nil")
		}
	})

	t.Run("an invalid email", func(t *testing.T) {
		_, err := uc.Login("client2@mail.com", password)
		if err == nil {
			t.Error("expected an error, but got a nil")
		}

		if err != nil && err.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, err.Type)
		}
	})

	t.Run("an invalid password", func(t *testing.T) {
		_, err := uc.Login(email, "123456")
		if err == nil {
			t.Error("expected an error, but got a nil")
		}

		if err != nil && err.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, err.Type)
		}
	})
}
