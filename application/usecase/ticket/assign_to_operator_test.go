package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"testing"
)

func TestAssignToOperator(t *testing.T) {
	uc := newTicketUC()
	operator, _ := domain.UserRegister("operator", "operator@mail.com", "password")
	operator2, _ := domain.UserRegister("operator2", "operator2@mail.com", "password")
	operator.Profile = domain.PROFILE_OPERATOR
	operator2.Profile = domain.PROFILE_OPERATOR
	_, _ = uc.UserUseCase.UserRepository.Insert(operator)
	_, _ = uc.UserUseCase.UserRepository.Insert(operator2)
	operatorResponse, _ := uc.UserUseCase.Login("operator@mail.com", "password")
	client, _ := domain.UserRegister("client", "client@mail.com", "password")
	_, _ = uc.UserUseCase.UserRepository.Insert(client)
	clientResponse, _ := uc.UserUseCase.Login("client@mail.com", "password")

	t.Run("assign a valid ticket", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)
		got, _ := uc.AssignToOperator(ticket.ID, operator2.ID, operatorResponse.Token)

		if got == nil {
			t.Error("expected a ticket, but got a nil")
		}
	})

	t.Run("a no-existent ticket", func(t *testing.T) {
		_, myerr := uc.AssignToOperator("123", operator2.ID, operatorResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.REGISTER_NOT_FOUND {
			t.Errorf("expected %q, but got %q", myerrors.REGISTER_NOT_FOUND, myerr.Type)
		}
	})

	t.Run("an invalid token", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.AssignToOperator(ticket.ID, operator2.ID, "123")

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.UNAUTHORIZED {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})

	t.Run("assign a ticket with a client", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.AssignToOperator(ticket.ID, operator2.ID, clientResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.FORBIDDEN {
			t.Errorf("expected %q, but got %q", myerrors.FORBIDDEN, myerr.Type)
		}
	})

	t.Run("assign a ticket to a client", func(t *testing.T) {
		ticket, _ := uc.Open("title", "description", clientResponse.Token)

		_, myerr := uc.AssignToOperator(ticket.ID, client.ID, operatorResponse.Token)

		if myerr == nil {
			t.Error("expected an error, but got a nil")
		}
		if myerr.Type != myerrors.DOMAIN {
			t.Errorf("expected %q, but got %q", myerrors.DOMAIN, myerr.Type)
		}
	})
}
