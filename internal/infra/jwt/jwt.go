package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
)

type Payload struct {
	ID        string
	Name      string
	Profile   string
	ExpiresAt time.Time
}

func (payload Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return security.ErrUnauthorized
	}
	return nil
}

type Authenticator struct {
	secret   string
	duration time.Duration
}

func NewAuthenticator(secret string, duration time.Duration) *Authenticator {
	return &Authenticator{
		secret:   secret,
		duration: duration,
	}
}

func (a *Authenticator) Generate(user security.User) (string, error) {
	payload := Payload{
		ID:        user.ID,
		Name:      user.Name,
		Profile:   user.Profile,
		ExpiresAt: time.Now().Add(a.duration),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := at.SignedString([]byte(a.secret))
	if err != nil {
		return "", security.ErrUnauthorized
	}

	return token, nil
}

func (a *Authenticator) Validate(token string) (*security.User, error) {
	payload := Payload{}
	_, err := jwt.ParseWithClaims(token, &payload, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, security.ErrUnauthorized
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, security.ErrUnauthorized
	}

	return &security.User{
		ID:      payload.ID,
		Name:    payload.Name,
		Profile: payload.Profile,
	}, nil
}
