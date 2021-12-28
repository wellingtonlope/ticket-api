package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestGetById(t *testing.T) {
	t.Run("get a ticket with an operator", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
		ticket, _ := uc.Open("title", "description", operatorResponse.Token)

		got, _ := uc.GetById(ticket.ID, operatorResponse.Token)
		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("get a ticket with a client", func(t *testing.T) {
		uc := newTicketUC()
		clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		got, _ := uc.GetById(ticket.ID, clientResponse.Token)
		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("get a ticket with a operator from another user", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
		clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		got, _ := uc.GetById(ticket.ID, operatorResponse.Token)
		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("get a ticket with a client from another user", func(t *testing.T) {
		uc := newTicketUC()
		clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")
		clientResponse2, _ := uc.UserUseCase.Register("client", "client2@mail.com", "password")
		ticket, _ := uc.Open("title", "description", clientResponse2.Token)

		_, myerr := uc.GetById(ticket.ID, clientResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but got %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})

	t.Run("get a ticket with an invalid id", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")

		_, myerr := uc.GetById("123", operatorResponse.Token)
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, myerr.Type)
		}
	})

	t.Run("get a ticket with an invalid token", func(t *testing.T) {
		uc := newTicketUC()
		operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
		operator.Profile = domain.PROFILE_OPERATOR
		_, _ = uc.UserUseCase.UserRepository.Insert(operator)
		operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
		ticket, _ := uc.Open("title", "description", operatorResponse.Token)

		_, myerr := uc.GetById(ticket.ID, "123")
		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})
}
