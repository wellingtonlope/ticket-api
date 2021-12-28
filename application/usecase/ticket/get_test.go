package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestGet(t *testing.T) {
	uc := newTicketUC()
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	_, _ = uc.UserUseCase.UserRepository.Insert(operator)
	operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
	clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")

	t.Run("get a valid ticket", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)
		got, _ := uc.Get(ticket.ID, operatorResponse.Token)

		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("a no-existent ticket", func(t *testing.T) {
		_, myerr := uc.Get("123", operatorResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, myerr.Type)
		}
	})

	t.Run("an invalid token", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.Get(ticket.ID, "123")

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})

	t.Run("get a ticket with a client", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.Get(ticket.ID, clientResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but got %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})
}
