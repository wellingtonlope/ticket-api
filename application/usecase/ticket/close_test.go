package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestClose(t *testing.T) {
	uc := newTicketUC()
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	_, _ = uc.UserUseCase.UserRepository.Insert(operator)
	operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
	clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")

	t.Run("close a valid ticket", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)
		_, _ = uc.Get(ticket.ID, operatorResponse.Token)

		got, _ := uc.Close(ticket.ID, "solution", operatorResponse.Token)
		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("close an unpicked ticket", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.Close(ticket.ID, "solution", operatorResponse.Token)
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.DOMAIN {
			t.Errorf("expected %q, but got %q", myerrors.DOMAIN, myerr.Type)
		}
	})

	t.Run("close a no-existent ticket", func(t *testing.T) {
		_, myerr := uc.Close("123", "solution", operatorResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, myerr.Type)
		}
	})

	t.Run("close with an invalid token", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)
		_, _ = uc.Get(ticket.ID, operatorResponse.Token)

		_, myerr := uc.Close(ticket.ID, "solution", "123")

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})

	t.Run("close with a client", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)
		_, _ = uc.Get(ticket.ID, operatorResponse.Token)

		_, myerr := uc.Close(ticket.ID, "solution", clientResponse.Token)
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but got %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})
}
