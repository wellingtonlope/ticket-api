package ticket

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	uc := newTicketUC()
	clientResponse, _ := uc.UserUseCase.Register("client", "client@mail.com", "password")

	t.Run("a valid ticket", func(t *testing.T) {
		got, _ := uc.Open("title", "description", clientResponse.Token)

		if got == nil {
			t.Errorf("expected a ticket, but got a nil")
		}
	})

	t.Run("an invalid ticket without title", func(t *testing.T) {
		_, myerr := uc.Open("", "description", clientResponse.Token)

		if myerr == nil {
			t.Errorf("expected an error, but got a nil")
		}
		if myerrors.DOMAIN != myerr.Type {
			t.Errorf("expected %q, but got %q", myerrors.DOMAIN, myerr.Type)
		}
	})

	t.Run("an invalid token", func(t *testing.T) {
		ucExpired := newTicketUC()
		ucExpired.TokenUseCase.Duration = time.Minute * -1

		clientExpiredResponse, _ := ucExpired.UserUseCase.Register("client", "client@mail.com", "password")

		_, myerr := ucExpired.Open("title", "description", clientExpiredResponse.Token)

		if myerr == nil {
			t.Errorf("expected an error, but got a nil")
		}
		if myerrors.UNAUTHORIZED != myerr.Type {
			t.Errorf("expected %q, but got %q", myerrors.UNAUTHORIZED, myerr.Type)
		}
	})
}
